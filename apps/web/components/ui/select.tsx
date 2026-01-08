import { forwardRef, type SelectHTMLAttributes } from 'react';
import { ChevronDown } from 'lucide-react';
import styles from './select.module.css';

export interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  error?: string;
  hint?: string;
  options: Array<{ value: string; label: string }>;
  placeholder?: string;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  (
    { label, error, hint, options, placeholder, className = '', id, ...props },
    ref
  ) => {
    const selectId = id || `select-${Math.random().toString(36).slice(2)}`;

    return (
      <div className={`${styles.wrapper} ${className}`}>
        {label && (
          <label htmlFor={selectId} className={styles.label}>
            {label}
          </label>
        )}
        <div className={`${styles.selectWrapper} ${error ? styles.hasError : ''}`}>
          <select ref={ref} id={selectId} className={styles.select} {...props}>
            {placeholder && (
              <option value="" disabled>
                {placeholder}
              </option>
            )}
            {options.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          <ChevronDown size={18} className={styles.chevron} />
        </div>
        {error && <span className={styles.error}>{error}</span>}
        {!error && hint && <span className={styles.hint}>{hint}</span>}
      </div>
    );
  }
);

Select.displayName = 'Select';
