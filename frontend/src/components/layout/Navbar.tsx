import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';

export default function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const { user, logout } = useAuth();
  const location = useLocation();

  const navLinks = [
    { to: '/', label: 'Home' },
    { to: 'https://starkdb.org', label: 'STARKDB', external: true },
    { to: 'https://starkdb.org/docs', label: 'Docs', external: true },
    { to: 'https://lab.starkdb.org', label: 'Lab', external: true },
    { to: 'https://github.com/abdullahtnz/starkdb', label: 'GitHub', external: true },
  ];

  return (
    <>
      <nav className="fixed top-0 left-0 right-0 h-16 bg-[#0a0a0f]/85 backdrop-blur-xl border-b border-stark-border z-50 flex items-center">
        <div className="nav-container">
          <Link to="/" className="flex items-center gap-2.5 font-display font-extrabold text-lg text-stark-text tracking-[-0.02em]">
            <img src="/starklogo.png" alt="STARKDB" className="h-8 w-8" />
            STARK<span className="text-stark-accent">DB</span>
          </Link>

          <ul className="hidden md:flex gap-7 items-center list-none">
            {navLinks.map((link) =>
              link.external ? (
                <li key={link.to}>
                  <a
                    href={link.to}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-mono text-xs tracking-[0.06em] text-stark-muted uppercase relative py-1 hover:text-stark-accent transition-colors"
                  >
                    {link.label}
                  </a>
                </li>
              ) : (
                <li key={link.to}>
                  <Link
                    to={link.to}
                    className={`font-mono text-xs tracking-[0.06em] uppercase relative py-1 transition-colors ${
                      location.pathname === link.to ? 'text-stark-accent' : 'text-stark-muted hover:text-stark-accent'
                    }`}
                  >
                    {link.label}
                  </Link>
                </li>
              )
            )}
          </ul>

          <div className="hidden md:flex items-center gap-4">
            {user ? (
              <>
                <Link
                  to="/ask"
                  className="btn-primary text-xs py-2 px-4 tracking-[0.08em]"
                >
                  Ask Question
                </Link>
                <div className="relative group">
                  <button className="flex items-center gap-2 text-stark-text hover:text-stark-accent transition-colors">
                    <span className="font-mono text-xs tracking-[0.06em] uppercase">{user.username}</span>
                    <div className="w-7 h-7 rounded bg-stark-accent/10 border border-stark-accent/20 flex items-center justify-center text-xs font-mono text-stark-accent">
                      {user.username[0].toUpperCase()}
                    </div>
                  </button>
                  <div className="absolute right-0 top-full mt-2 w-48 bg-stark-surface border border-stark-border opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
                    <Link to="/profile" className="block px-4 py-3 text-sm text-stark-muted hover:text-stark-text hover:bg-stark-surface-2 transition-colors font-mono">Profile</Link>
                    {user.is_admin && (
                      <Link to="/admin" className="block px-4 py-3 text-sm text-stark-muted hover:text-stark-text hover:bg-stark-surface-2 transition-colors font-mono">Admin</Link>
                    )}
                    <button onClick={logout} className="w-full text-left px-4 py-3 text-sm text-stark-accent-3 hover:bg-stark-surface-2 transition-colors font-mono">Logout</button>
                  </div>
                </div>
              </>
            ) : (
              <>
                <Link to="/login" className="font-mono text-xs tracking-[0.06em] text-stark-muted hover:text-stark-accent transition-colors uppercase">Login</Link>
                <Link to="/signup" className="btn-primary text-xs py-2 px-4 tracking-[0.08em]">Sign Up</Link>
              </>
            )}
          </div>

          <button
            className="md:hidden bg-transparent border-none cursor-pointer w-9 h-9 relative z-50"
            onClick={() => setMobileOpen(!mobileOpen)}
            aria-label="Toggle menu"
          >
            <span className={`block w-[22px] h-[2px] bg-stark-text absolute left-[7px] transition-all rounded-sm ${mobileOpen ? 'top-[17px] rotate-45' : 'top-[11px]'}`} />
            <span className={`block w-[22px] h-[2px] bg-stark-text absolute left-[7px] top-[17px] transition-all ${mobileOpen ? 'opacity-0' : ''}`} />
            <span className={`block w-[22px] h-[2px] bg-stark-text absolute left-[7px] transition-all rounded-sm ${mobileOpen ? 'top-[17px] -rotate-45' : 'top-[23px]'}`} />
          </button>
        </div>
      </nav>

      {mobileOpen && (
        <div className="fixed inset-0 bg-[#0a0a0f]/97 backdrop-blur-xl z-40 flex flex-col items-center justify-center gap-6 md:hidden">
          {navLinks.map((link) =>
            link.external ? (
              <a key={link.to} href={link.to} target="_blank" rel="noopener noreferrer" className="font-display text-2xl font-bold text-stark-muted hover:text-stark-accent transition-colors">{link.label}</a>
            ) : (
              <Link key={link.to} to={link.to} onClick={() => setMobileOpen(false)} className={`font-display text-2xl font-bold transition-colors ${location.pathname === link.to ? 'text-stark-accent' : 'text-stark-muted hover:text-stark-accent'}`}>{link.label}</Link>
            )
          )}
          {user ? (
            <>
              <Link to="/ask" onClick={() => setMobileOpen(false)} className="btn-primary text-sm py-3 px-7 mt-4">Ask Question</Link>
              <Link to="/profile" onClick={() => setMobileOpen(false)} className="font-display text-2xl font-bold text-stark-muted hover:text-stark-accent transition-colors">Profile</Link>
              {user.is_admin && <Link to="/admin" onClick={() => setMobileOpen(false)} className="font-display text-2xl font-bold text-stark-muted hover:text-stark-accent transition-colors">Admin</Link>}
              <button onClick={() => { logout(); setMobileOpen(false); }} className="font-display text-2xl font-bold text-stark-accent-3 hover:text-stark-accent transition-colors">Logout</button>
            </>
          ) : (
            <>
              <Link to="/login" onClick={() => setMobileOpen(false)} className="font-display text-2xl font-bold text-stark-muted hover:text-stark-accent transition-colors">Login</Link>
              <Link to="/signup" onClick={() => setMobileOpen(false)} className="btn-primary text-sm py-3 px-7 mt-4">Sign Up</Link>
            </>
          )}
        </div>
      )}
    </>
  );
}
