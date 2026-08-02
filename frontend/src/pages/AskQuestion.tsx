import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api'
import { useAuth } from '../contexts/AuthContext'
import ProtectedRoute from '../components/auth/ProtectedRoute'
import type { Tag } from '../types'

export default function AskQuestion() {
  return (
    <ProtectedRoute>
      <AskQuestionForm />
    </ProtectedRoute>
  )
}

function AskQuestionForm() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [tags, setTags] = useState<Tag[]>([])
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const [uploadedFiles, setUploadedFiles] = useState<{ id: string; file_name: string; file_path: string }[]>([])
  const [uploading, setUploading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    api.get<Tag[]>('/tags').then(data => setTags(data || [])).catch(() => {})
  }, [])

  const toggleTag = (tagId: string) => {
    setSelectedTags(prev => prev.includes(tagId) ? prev.filter(t => t !== tagId) : [...prev, tagId])
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    try {
      const result = await api.upload(file)
      setUploadedFiles(prev => [...prev, result])
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setUploading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (title.length < 5) { setError('Title must be at least 5 characters'); return }
    if (!content.trim()) { setError('Content cannot be empty'); return }
    setSubmitting(true)
    setError('')
    try {
      const post = await api.post<{ id: string }>('/posts', {
        title,
        content,
        tag_ids: selectedTags,
      })
      navigate(`/posts/${post.id}`)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create post')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen pt-20 pb-12">
      <div className="max-w-3xl mx-auto px-4 sm:px-8">
        <div className="card">
          <div className="section-tag">New Discussion</div>
          <h1 className="text-3xl font-display font-extrabold text-stark-text tracking-[-0.03em] mb-8">Ask a Question</h1>

          {error && (
            <div className="bg-red-500/10 border border-red-500/30 text-red-400 font-mono text-sm p-3 mb-6">{error}</div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Title</label>
              <input
                type="text"
                value={title}
                onChange={e => setTitle(e.target.value)}
                required
                placeholder="e.g. How to optimize B-Tree lookups in STARKDB?"
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
              />
            </div>

            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Content</label>
              <textarea
                value={content}
                onChange={e => setContent(e.target.value)}
                required
                rows={10}
                placeholder="Describe your question or solution in detail. You can use HTML tags for formatting."
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50 resize-y"
              />
            </div>

            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Tags</label>
              <div className="flex flex-wrap gap-2">
                {tags.map(tag => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => toggleTag(tag.id)}
                    className={`font-mono text-xs px-3 py-1.5 border transition-colors ${
                      selectedTags.includes(tag.id)
                        ? 'bg-stark-accent/10 border-stark-accent/30 text-stark-accent'
                        : 'bg-stark-surface border-stark-border text-stark-muted hover:border-stark-accent/30 hover:text-stark-accent'
                    }`}
                  >
                    {tag.name}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Attachments</label>
              <div className="flex items-center gap-3">
                <input
                  type="file"
                  ref={fileInputRef}
                  onChange={handleFileUpload}
                  accept="image/*,.pdf,.txt,.csv,.zip"
                  className="font-mono text-xs text-stark-muted file:mr-3 file:py-2 file:px-4 file:bg-stark-surface file:border file:border-stark-border file:text-stark-accent file:font-mono file:text-xs file:cursor-pointer file:hover:border-stark-accent file:transition-colors"
                />
                {uploading && <span className="font-mono text-xs text-stark-muted">Uploading...</span>}
              </div>
              {uploadedFiles.length > 0 && (
                <div className="flex flex-wrap gap-2 mt-3">
                  {uploadedFiles.map(file => (
                    <span key={file.id} className="font-mono text-xs px-2 py-1 bg-stark-surface-2 border border-stark-border text-stark-accent">
                      {file.file_name}
                    </span>
                  ))}
                </div>
              )}
            </div>

            <button type="submit" disabled={submitting} className="btn-primary">
              {submitting ? 'Posting...' : 'Post Question'}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}
