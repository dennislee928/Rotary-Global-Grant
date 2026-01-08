import { type HTMLAttributes } from 'react';
import styles from './badge.module.css';

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'severity';
  severity?: 'S0' | 'S1' | 'S2' | 'S3' | 'S4';
  size?: 'sm' | 'md';
}

export function Badge({
  variant = 'default',
  severity,
  size = 'md',
  className = '',
  children,
  ...props
}: BadgeProps) {
  const variantClass = severity ? styles[`severity-${severity}`] : styles[variant];
  
  return (
    <span
      className={`${styles.badge} ${variantClass} ${styles[size]} ${className}`}
      {...props}
    >
      {children}
    </span>
  );
}
