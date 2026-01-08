'use client';

import { useEffect, useState } from 'react';
import { AlertTriangle, Bell, CheckCircle, XCircle, Clock, Send } from 'lucide-react';
import { getAlerts } from '@/lib/api';
import type { Alert } from '@/lib/types';
import { Card, CardContent, Badge, Button } from '@/components/ui';
import styles from './page.module.css';

const STATUS_CONFIG = {
  draft: { icon: Clock, color: 'default' as const, label: 'Draft' },
  approved: { icon: CheckCircle, color: 'info' as const, label: 'Approved' },
  published: { icon: Send, color: 'success' as const, label: 'Published' },
  withdrawn: { icon: XCircle, color: 'error' as const, label: 'Withdrawn' },
};

export default function AlertsPage() {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<string>('');

  useEffect(() => {
    async function loadAlerts() {
      try {
        const result = await getAlerts({
          pageSize: 50,
          status: filter || undefined,
        });
        setAlerts(result.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load alerts');
      } finally {
        setLoading(false);
      }
    }
    loadAlerts();
  }, [filter]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>Alert Management</h1>
          <p className={styles.subtitle}>Create and manage CAP-ready alerts</p>
        </div>
        <Button leftIcon={<Bell size={18} />}>
          Create Alert
        </Button>
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
        {Object.entries(STATUS_CONFIG).map(([status, config]) => (
          <Button
            key={status}
            variant={filter === status ? 'primary' : 'ghost'}
            size="sm"
            onClick={() => setFilter(status)}
          >
            {config.label}
          </Button>
        ))}
      </div>

      {/* Alerts List */}
      {loading ? (
        <Card>
          <CardContent>
            <p className={styles.loadingState}>Loading alerts...</p>
          </CardContent>
        </Card>
      ) : error ? (
        <Card>
          <CardContent>
            <div className={styles.errorState}>
              <AlertTriangle size={24} />
              <p>{error}</p>
            </div>
          </CardContent>
        </Card>
      ) : alerts.length === 0 ? (
        <Card>
          <CardContent>
            <p className={styles.emptyState}>No alerts found</p>
          </CardContent>
        </Card>
      ) : (
        <div className={styles.alertsGrid}>
          {alerts.map((alert) => {
            const statusConfig = STATUS_CONFIG[alert.status];
            const StatusIcon = statusConfig.icon;
            return (
              <Card key={alert.id} className={styles.alertCard}>
                <CardContent>
                  <div className={styles.alertHeader}>
                    <Badge variant={statusConfig.color}>
                      <StatusIcon size={12} />
                      {statusConfig.label}
                    </Badge>
                    <span className={styles.alertTime}>
                      {formatDate(alert.createdAt)}
                    </span>
                  </div>

                  <h3 className={styles.alertEvent}>{alert.event}</h3>

                  <div className={styles.alertMeta}>
                    <Badge size="sm" variant={
                      alert.severity === 'Extreme' || alert.severity === 'Severe' 
                        ? 'error' 
                        : alert.severity === 'Moderate' 
                          ? 'warning' 
                          : 'info'
                    }>
                      {alert.severity}
                    </Badge>
                    <Badge size="sm">{alert.urgency}</Badge>
                    <Badge size="sm">{alert.certainty}</Badge>
                  </div>

                  <div className={styles.alertDetails}>
                    <div className={styles.alertField}>
                      <span className={styles.fieldLabel}>Area</span>
                      <span className={styles.fieldValue}>{alert.area}</span>
                    </div>
                    <div className={styles.alertField}>
                      <span className={styles.fieldLabel}>Instructions</span>
                      <p className={styles.fieldValue}>{alert.instruction}</p>
                    </div>
                    {alert.publicMessage && (
                      <div className={styles.alertField}>
                        <span className={styles.fieldLabel}>Public Message</span>
                        <p className={styles.fieldValue}>{alert.publicMessage}</p>
                      </div>
                    )}
                  </div>

                  {alert.status === 'draft' && (
                    <div className={styles.alertActions}>
                      <Button size="sm" variant="outline">
                        Edit
                      </Button>
                      <Button size="sm">
                        Approve
                      </Button>
                    </div>
                  )}
                  {alert.status === 'approved' && (
                    <div className={styles.alertActions}>
                      <Button size="sm" variant="outline">
                        Edit
                      </Button>
                      <Button size="sm" leftIcon={<Send size={14} />}>
                        Publish
                      </Button>
                    </div>
                  )}
                  {alert.publishedAt && (
                    <div className={styles.publishedInfo}>
                      Published: {formatDate(alert.publishedAt)}
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
