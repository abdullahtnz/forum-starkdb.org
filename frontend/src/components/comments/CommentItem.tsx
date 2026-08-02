import { useState } from 'react'
import { api } from '../../api'
import { useAuth } from '../../contexts/AuthContext'
import type { Comment as CommentType } from '../../types'

function timeAgo(date: string) {
  const seconds = Math.floor((Date.now() - new Date(date).getTime()) / 1000)
  if (seconds < 60) return 'just now'
  const mins = Math.floor(seconds / 60)
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}d ago`
  const months = Math.floor(days / 30)
  if (months < 12) return `${months}mo ago`
  return `${Math.floor(months / 12)}y ago`
}

export default function CommentItem({
  comment,
  postId,
  onRefresh,
}: {
  comment: CommentType
  postId: string
  onRefresh: () => void
}) {
  const { user } = useAuth()
  const [showReply, setShowReply] = useState(false)
  const [replyText, setReplyText] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [showReplies, setShowReplies] = useState(true)
  const [liked, setLiked] = useState(comment.liked)
  const [likeCount, setLikeCount] = useState(comment.like_count)
  const replies = comment.replies || []

  const handleLike = async () => {
    if (!user) return
    try {
      const res = await api.post<{ liked: boolean }>(`/comments/${comment.id}/like`, {})
      setLiked(res.liked)
      setLikeCount(c => res.liked ? c + 1 : c - 1)
    } catch {}
  }

  const handleReply = async () => {
    if (!replyText.trim() || !user) return
    setSubmitting(true)
    try {
      await api.post(`/comments/${comment.id}/reply`, {
        content: replyText,
        post_id: postId,
      })
      setReplyText('')
      setShowReply(false)
      onRefresh()
    } catch {} finally {
      setSubmitting(false)
    }
  }

  if (comment.is_deleted && replies.length === 0) {
    return (
      <div className="pl-4 border-l-2 border-stark-border py-2">
        <p className="font-mono text-xs text-stark-muted italic">[deleted]</p>
      </div>
    )
  }

  return (
    <div className="pl-4 border-l-2 border-stark-border">
      <div className="py-3">
        <div className="flex items-center gap-2 mb-1.5">
          <div className="w-5 h-5 rounded bg-stark-accent/10 border border-stark-accent/20 flex items-center justify-center text-[0.55rem] font-mono text-stark-accent font-bold">
            {comment.username[0]?.toUpperCase()}
          </div>
          <span className="font-mono text-xs text-stark-accent">{comment.username}</span>
          <span className="font-mono text-[0.6rem] text-stark-muted">· {timeAgo(comment.created_at)}</span>
          {comment.is_owner && (
            <span className="font-mono text-[0.55rem] bg-stark-accent/10 px-1.5 py-0.5 text-stark-accent border border-stark-accent/20">OP</span>
          )}
        </div>
        <p className="font-mono text-xs text-stark-muted leading-relaxed ml-7" dangerouslySetInnerHTML={{ __html: comment.content }} />

        <div className="flex items-center gap-4 ml-7 mt-1.5">
          <button
            onClick={handleLike}
            disabled={!user}
            className={`font-mono text-[0.6rem] uppercase tracking-[0.1em] transition-colors ${
              liked ? 'text-stark-accent' : 'text-stark-muted hover:text-stark-accent'
            } disabled:cursor-not-allowed`}
          >
            {liked ? '♥' : '♡'} {likeCount}
          </button>
          {user && (
            <button
              onClick={() => setShowReply(!showReply)}
              className="font-mono text-[0.6rem] uppercase tracking-[0.1em] text-stark-muted hover:text-stark-accent transition-colors"
            >
              Reply
            </button>
          )}
        </div>

        {showReply && user && (
          <div className="ml-7 mt-3 flex gap-2">
            <input
              type="text"
              value={replyText}
              onChange={e => setReplyText(e.target.value)}
              placeholder="Write a reply..."
              className="flex-1 bg-stark-bg border border-stark-border text-stark-text font-mono text-xs p-2 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
            />
            <button onClick={handleReply} disabled={submitting || !replyText.trim()} className="btn-primary text-xs py-2 px-4">
              {submitting ? '...' : 'Reply'}
            </button>
            <button onClick={() => setShowReply(false)} className="btn-ghost text-xs py-2 px-3">Cancel</button>
          </div>
        )}
      </div>

      {replies.length > 0 && (
        <div>
          <button
            onClick={() => setShowReplies(!showReplies)}
            className="font-mono text-[0.6rem] text-stark-muted hover:text-stark-accent transition-colors mb-2 ml-7"
          >
            {showReplies ? '▾ Hide' : '▸ Show'} {replies.length} {replies.length === 1 ? 'reply' : 'replies'}
          </button>
          {showReplies && (
            <div className="ml-2">
              {replies.map(reply => (
                <CommentItem key={reply.id} comment={reply} postId={postId} onRefresh={onRefresh} />
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
