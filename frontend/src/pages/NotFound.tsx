import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center px-4">
      <div className="bg-stark-surface border border-stark-border p-12 max-w-lg w-full text-center">
        <div className="text-8xl font-display font-extrabold text-stark-accent animate-pulse mb-2">404</div>
        <h1 className="text-2xl font-display font-bold text-stark-text mb-4">Page Not Found</h1>
        <p className="font-mono text-sm text-stark-muted mb-8">
          The page you are looking for does not exist or has been moved.
        </p>
        <Link to="/" className="btn-primary justify-center">
          Back to Home
        </Link>
      </div>
    </div>
  )
}
