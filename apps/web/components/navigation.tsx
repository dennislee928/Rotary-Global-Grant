'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { 
  Shield, 
  FileText, 
  AlertTriangle, 
  GraduationCap, 
  BarChart3,
  Menu,
  X 
} from 'lucide-react';
import { ThemeToggle } from './theme-toggle';
import { useState } from 'react';
import styles from './navigation.module.css';

const navItems = [
  { href: '/', label: 'Dashboard', icon: BarChart3 },
  { href: '/report', label: 'Report', icon: FileText },
  { href: '/triage', label: 'Triage', icon: Shield },
  { href: '/alerts', label: 'Alerts', icon: AlertTriangle },
  { href: '/training', label: 'Training', icon: GraduationCap },
  { href: '/kpi', label: 'KPIs', icon: BarChart3 },
];

export function Navigation() {
  const pathname = usePathname();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className={styles.header}>
      <div className={styles.container}>
        <Link href="/" className={styles.logo}>
          <Shield size={28} />
          <span className={styles.logoText}>The Hive</span>
        </Link>

        <nav className={`${styles.nav} ${mobileMenuOpen ? styles.navOpen : ''}`}>
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`${styles.navLink} ${isActive ? styles.navLinkActive : ''}`}
                onClick={() => setMobileMenuOpen(false)}
              >
                <Icon size={18} />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>

        <div className={styles.actions}>
          <ThemeToggle />
          <button
            className={styles.mobileMenuButton}
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            aria-label="Toggle menu"
          >
            {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>
      </div>
    </header>
  );
}
