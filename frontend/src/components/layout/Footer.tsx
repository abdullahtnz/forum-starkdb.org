import { Link } from 'react-router-dom';

export default function Footer() {
  return (
    <footer className="border-t border-stark-border py-8 text-center max-w-7xl mx-auto px-8">
      <p className="font-mono text-xs text-stark-border mb-3 tracking-[0.1em]">STARK · Open Source · MIT License</p>
      <div className="flex justify-center gap-6">
        <a href="https://github.com/abdullahtnz/starkdb" target="_blank" rel="noopener noreferrer" className="font-mono text-xs text-stark-muted tracking-[0.06em] hover:text-stark-accent transition-colors">GitHub</a>
        <a href="https://starkdb.org/docs" target="_blank" rel="noopener noreferrer" className="font-mono text-xs text-stark-muted tracking-[0.06em] hover:text-stark-accent transition-colors">Docs</a>
        <Link to="/privacy" className="font-mono text-xs text-stark-muted tracking-[0.06em] hover:text-stark-accent transition-colors">Privacy Policy</Link>
        <Link to="/terms" className="font-mono text-xs text-stark-muted tracking-[0.06em] hover:text-stark-accent transition-colors">Terms of Use</Link>
      </div>
    </footer>
  );
}
