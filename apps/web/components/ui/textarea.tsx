import { forwardRef, type TextareaHTMLAttributes } from 'react';
import styles from './textarea.module.css';

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
  hint?: string;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ label, error, hint, className = '', id, ...props }, ref) => {
    const textareaId = id || `textarea-${Math.random().toString(36).slice(2)}`;

    return (
      <div className={`${styles.wrapper} ${className}`}>
        {label && (
          <label htmlFor={textareaId} className={styles.label}>
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          id={textareaId}
          className={`${styles.textarea} ${error ? styles.hasError : ''}`}
          {...props}
        />
        {error && <span className={styles.error}>{error}</span>}
        {!error && hint && <span className={styles.hint}>{hint}</span>}
      </div>
    );
  }
);

Textarea.displayName = 'Textarea';
