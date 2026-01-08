'use client';

import { useEffect, useState } from 'react';
import { Shield, Clock, AlertCircle, CheckCircle, XCircle, HelpCircle, ArrowUpCircle } from 'lucide-react';
import { getReports } from '@/lib/api';
import type { Report } from '@/lib/types';
import { Card, CardContent, CardHeader, CardTitle, Badge, Button } from '@/components/ui';
import styles from './page.module.css';

const STATUS_ICONS = {
  submitted: Clock,
  under_review: AlertCircle,
  triaged: CheckCircle,
  escalated: ArrowUpCircle,
  closed: CheckCircle,
  spam: XCircle,
};

const STATUS_COLORS: Record<string, 'default' | 'info' | 'success' | 'warning' | 'error'> = {
  submitted: 'info',
  under_review: 'warning',
  triaged: 'success',
  escalated: 'error',
  closed: 'default',
  spam: 'error',
};

export default function TriagePage() {
  const [reports, setReports] = useState<Report[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<string>('');

  useEffect(() => {
    async function loadReports() {
      try {
        const result = await getReports({
          pageSize: 50,
          status: filter || undefined,
        });
        setReports(result.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load reports');
      } finally {
        setLoading(false);
      }
    }
    loadReports();
  }, [filter]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>Triage Console</h1>
          <p className={styles.subtitle}>Review and process community reports</p>
        </div>
      </header>

      {/* Filters */}
      <div className={styles.filters}>
        <Button
          variant={filter === '' ? 'primary' : 'ghost'}
          size="sm"
          onClick={() => setFilter('')}
        >
          All
        </Button>
        <Button
          variant={filter === 'submitted' ? 'primary' : 'ghost'}
          size="sm"
          onClick={() => setFilter('submitted')}
        >
          Pending
        </Button>
        <Button
          variant={filter === 'under_review' ? 'primary' : 'ghost'}
          size="sm"
          onClick={() => setFilter('under_review')}
        >
          Under Review
        </Button>
        <Button
          variant={filter === 'triaged' ? 'primary' : 'ghost'}
          size="sm"
          onClick={() => setFilter('triaged')}
        >
          Triaged
        </Button>
        <Button
          variant={filter === 'escalated' ? 'primary' : 'ghost'}
          size="sm"
          onClick={() => setFilter('escalated')}
        >
          Escalated
        </Button>
      </div>

      {/* Reports List */}
      {loading ? (
        <Card>
          <CardContent>
            <p className={styles.loadingState}>Loading reports...</p>
          </CardContent>
        </Card>
      ) : error ? (
        <Card>
          <CardContent>
            <div className={styles.errorState}>
              <AlertCircle size={24} />
              <p>{error}</p>
              <p className={styles.errorHint}>
                Make sure you are logged in and have triage permissions.
              </p>
            </div>
          </CardContent>
        </Card>
      ) : reports.length === 0 ? (
        <Card>
          <CardContent>
            <p className={styles.emptyState}>No reports found</p>
          </CardContent>
        </Card>
      ) : (
        <div className={styles.reportsList}>
          {reports.map((report) => {
            const StatusIcon = STATUS_ICONS[report.status] || HelpCircle;
            return (
              <Card key={report.id} className={styles.reportCard}>
                <CardContent>
                  <div className={styles.reportHeader}>
                    <div className={styles.reportMeta}>
                      <Badge>{report.category.replace(/_/g, ' ')}</Badge>
                      {report.severitySuggested && (
                        <Badge severity={report.severitySuggested}>
                          {report.severitySuggested}
                        </Badge>
                      )}
                      <Badge variant={STATUS_COLORS[report.status]}>
                        <StatusIcon size={12} />
                        {report.status}
                      </Badge>
                    </div>
                    <span className={styles.reportTime}>
                      {formatDate(report.createdAt)}
                    </span>
                  </div>
                  
                  <div className={styles.reportBody}>
                    <div className={styles.reportField}>
                      <span className={styles.fieldLabel}>Location</span>
                      <span className={styles.fieldValue}>{report.areaHint}</span>
                    </div>
                    {report.timeWindow && (
                      <div className={styles.reportField}>
                        <span className={styles.fieldLabel}>Time</span>
                        <span className={styles.fieldValue}>{report.timeWindow}</span>
                      </div>
                    )}
                    <div className={styles.reportField}>
                      <span className={styles.fieldLabel}>Description</span>
                      <p className={styles.fieldValue}>{report.description}</p>
                    </div>
                    {report.evidence && report.evidence.length > 0 && (
                      <div className={styles.reportField}>
                        <span className={styles.fieldLabel}>Evidence</span>
                        <span className={styles.fieldValue}>
                          {report.evidence.length} item(s)
                        </span>
                      </div>
                    )}
                  </div>

                  {(report.status === 'submitted' || report.status === 'under_review') && (
                    <div className={styles.reportActions}>
                      <Button size="sm" variant="outline">
                        Review
                      </Button>
                      <Button size="sm">
                        Triage
                      </Button>
                    </div>
                  )}
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
