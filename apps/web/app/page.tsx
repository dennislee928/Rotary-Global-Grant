import { 
  FileText, 
  AlertTriangle, 
  Users, 
  TrendingUp,
  ArrowRight
} from 'lucide-react';
import Link from 'next/link';
import { getDashboardStats, getHealth } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import styles from './page.module.css';

async function fetchData() {
  try {
    const [health, stats] = await Promise.all([
      getHealth().catch(() => null),
      getDashboardStats().catch(() => null),
    ]);
    return { health, stats };
  } catch {
    return { health: null, stats: null };
  }
}

export default async function Dashboard() {
  const { health, stats } = await fetchData();

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>Dashboard</h1>
          <p className={styles.subtitle}>
            Community Safety & Digital Resilience Overview
          </p>
        </div>
        <Link href="/report" className={styles.ctaButton}>
          Submit Report
          <ArrowRight size={18} />
        </Link>
      </header>

      {/* Status Banner */}
      <Card className={styles.statusBanner}>
        <CardContent>
          <div className={styles.statusGrid}>
            <div className={styles.statusItem}>
              <span className={styles.statusLabel}>API Status</span>
              <Badge variant={health?.status === 'ok' ? 'success' : 'error'}>
                {health?.status || 'Unknown'}
              </Badge>
            </div>
            <div className={styles.statusItem}>
              <span className={styles.statusLabel}>Database</span>
              <Badge variant={health?.database === 'connected' ? 'success' : 'warning'}>
                {health?.database || 'Unknown'}
              </Badge>
            </div>
            <div className={styles.statusItem}>
              <span className={styles.statusLabel}>Version</span>
              <span className={styles.statusValue}>{health?.version || 'N/A'}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Stats Grid */}
      <div className={styles.statsGrid}>
        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon} style={{ background: 'var(--color-info)' }}>
                <FileText size={24} />
              </div>
              <div className={styles.statInfo}>
                <span className={styles.statValue}>
                  {stats?.totalReports ?? '--'}
                </span>
                <span className={styles.statLabel}>Total Reports</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon} style={{ background: 'var(--color-success)' }}>
                <TrendingUp size={24} />
              </div>
              <div className={styles.statInfo}>
                <span className={styles.statValue}>
                  {stats?.reportsThisWeek ?? '--'}
                </span>
                <span className={styles.statLabel}>This Week</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon} style={{ background: 'var(--color-warning)' }}>
                <AlertTriangle size={24} />
              </div>
              <div className={styles.statInfo}>
                <span className={styles.statValue}>
                  {stats?.activeAlerts ?? '--'}
                </span>
                <span className={styles.statLabel}>Active Alerts</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon} style={{ background: 'var(--color-accent-primary)' }}>
                <Users size={24} />
              </div>
              <div className={styles.statInfo}>
                <span className={styles.statValue}>4</span>
                <span className={styles.statLabel}>Partner Orgs</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Recent Alerts */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>Recent Alerts</h2>
          <Link href="/alerts" className={styles.sectionLink}>
            View All <ArrowRight size={16} />
          </Link>
        </div>
        <Card>
          <CardContent>
            {stats?.recentAlerts && stats.recentAlerts.length > 0 ? (
              <div className={styles.alertList}>
                {stats.recentAlerts.map((alert) => (
                  <div key={alert.id} className={styles.alertItem}>
                    <div className={styles.alertInfo}>
                      <Badge 
                        variant={
                          alert.severity === 'Extreme' || alert.severity === 'Severe' 
                            ? 'error' 
                            : alert.severity === 'Moderate' 
                              ? 'warning' 
                              : 'info'
                        }
                        size="sm"
                      >
                        {alert.severity}
                      </Badge>
                      <span className={styles.alertEvent}>{alert.event}</span>
                    </div>
                    <span className={styles.alertArea}>{alert.area}</span>
                  </div>
                ))}
              </div>
            ) : (
              <p className={styles.emptyState}>No recent alerts</p>
            )}
          </CardContent>
        </Card>
      </section>

      {/* Category Breakdown */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>Reports by Category</h2>
        </div>
        <Card>
          <CardContent>
            {stats?.categoryBreakdown && stats.categoryBreakdown.length > 0 ? (
              <div className={styles.categoryGrid}>
                {stats.categoryBreakdown.map((cat) => (
                  <div key={cat.category} className={styles.categoryItem}>
                    <span className={styles.categoryName}>
                      {cat.category.replace(/_/g, ' ')}
                    </span>
                    <span className={styles.categoryCount}>{cat.count}</span>
                  </div>
                ))}
              </div>
            ) : (
              <p className={styles.emptyState}>No reports yet</p>
            )}
          </CardContent>
        </Card>
      </section>
    </div>
  );
}
