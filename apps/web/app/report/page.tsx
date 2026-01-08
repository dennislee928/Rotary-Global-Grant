'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Send, AlertCircle, CheckCircle } from 'lucide-react';
import { createReport } from '@/lib/api';
import type { CreateReportRequest, ReportCategory, SeverityLevel } from '@/lib/types';
import { Button, Card, CardContent, CardHeader, CardTitle, CardDescription, Input, Select, Textarea, Badge } from '@/components/ui';
import styles from './page.module.css';

const CATEGORIES: Array<{ value: ReportCategory; label: string }> = [
  { value: 'suspicious_item', label: 'Suspicious Item' },
  { value: 'suspicious_person', label: 'Suspicious Person' },
  { value: 'harassment_stalking', label: 'Harassment / Stalking' },
  { value: 'scam_phishing', label: 'Scam / Phishing' },
  { value: 'misinformation_panic', label: 'Misinformation / Panic' },
  { value: 'crowd_disorder', label: 'Crowd Disorder' },
  { value: 'infrastructure_hazard', label: 'Infrastructure Hazard' },
  { value: 'other', label: 'Other' },
];

const SEVERITIES: Array<{ value: SeverityLevel; label: string }> = [
  { value: 'S0', label: 'S0 - Informational' },
  { value: 'S1', label: 'S1 - Low Risk' },
  { value: 'S2', label: 'S2 - Moderate' },
  { value: 'S3', label: 'S3 - High' },
  { value: 'S4', label: 'S4 - Critical' },
];

export default function ReportPage() {
  const router = useRouter();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  
  const [formData, setFormData] = useState<CreateReportRequest>({
    category: 'other',
    areaHint: '',
    description: '',
    timeWindow: '',
    evidence: [],
    reporterContact: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      await createReport(formData);
      setSuccess(true);
      setTimeout(() => {
        router.push('/');
      }, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to submit report');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (success) {
    return (
      <div className={styles.container}>
        <Card className={styles.successCard}>
          <CardContent>
            <div className={styles.successContent}>
              <CheckCircle size={64} className={styles.successIcon} />
              <h2>Report Submitted</h2>
              <p>Thank you for your report. Our team will review it shortly.</p>
              <p className={styles.redirectText}>Redirecting to dashboard...</p>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>Submit a Report</h1>
        <p className={styles.subtitle}>
          Help keep our community safe by reporting suspicious activities or incidents
        </p>
      </header>

      <div className={styles.formContainer}>
        <Card>
          <CardHeader>
            <CardTitle>Incident Details</CardTitle>
            <CardDescription>
              Please provide as much detail as possible. All reports are reviewed by trained staff.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className={styles.form}>
              {error && (
                <div className={styles.errorBanner}>
                  <AlertCircle size={20} />
                  <span>{error}</span>
                </div>
              )}

              <Select
                label="Category *"
                options={CATEGORIES}
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value as ReportCategory })}
                required
              />

              <Select
                label="Suggested Severity"
                options={SEVERITIES}
                value={formData.severitySuggested || ''}
                onChange={(e) => setFormData({ ...formData, severitySuggested: e.target.value as SeverityLevel })}
                placeholder="Select severity level"
                hint="How urgent do you think this is?"
              />

              <Input
                label="Location / Area *"
                placeholder="e.g., Near the main entrance, Platform 3"
                value={formData.areaHint}
                onChange={(e) => setFormData({ ...formData, areaHint: e.target.value })}
                required
                hint="Approximate location (avoid exact addresses for privacy)"
              />

              <Input
                label="Time Window"
                placeholder="e.g., Today between 2-3 PM"
                value={formData.timeWindow}
                onChange={(e) => setFormData({ ...formData, timeWindow: e.target.value })}
                hint="When did this occur or when is it expected?"
              />

              <Textarea
                label="Description *"
                placeholder="Describe what you observed in detail..."
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                required
                rows={5}
                hint="Include relevant details but avoid personal identifiers"
              />

              <Input
                label="Evidence URLs"
                placeholder="https://..."
                value={formData.evidence?.join(', ') || ''}
                onChange={(e) => setFormData({ 
                  ...formData, 
                  evidence: e.target.value.split(',').map(s => s.trim()).filter(Boolean)
                })}
                hint="Comma-separated URLs to photos or documents (optional)"
              />

              <Input
                label="Contact (Optional)"
                placeholder="Email or phone for follow-up"
                value={formData.reporterContact}
                onChange={(e) => setFormData({ ...formData, reporterContact: e.target.value })}
                hint="Only used if we need to follow up"
              />

              <div className={styles.privacyNote}>
                <AlertCircle size={16} />
                <span>
                  Your report is handled according to our data minimization policy. 
                  We do not store unnecessary personal information.
                </span>
              </div>

              <div className={styles.formActions}>
                <Button type="button" variant="ghost" onClick={() => router.back()}>
                  Cancel
                </Button>
                <Button type="submit" isLoading={isSubmitting} leftIcon={<Send size={18} />}>
                  Submit Report
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
