import { Target, CheckCircle, XCircle, Clock, Users, AlertTriangle, BookOpen } from 'lucide-react';
import { getKPIMetrics } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import styles from './page.module.css';

async function fetchData() {
  try {
    const metrics = await getKPIMetrics();
    return { metrics };
  } catch {
    return { metrics: null };
  }
}

function KPICard({ 
  title, 
  current, 
  target, 
  unit = '', 
  inverse = false,
  format = (v: number) => v.toString()
}: { 
  title: string; 
  current: number; 
  target: number; 
  unit?: string;
  inverse?: boolean;
  format?: (v: number) => string;
}) {
  const isMet = inverse ? current <= target : current >= target;
  const progress = inverse 
    ? Math.max(0, Math.min(100, ((target - current) / target) * 100 + 100))
    : Math.min(100, (current / target) * 100);

  return (
    <Card>
      <CardContent>
        <div className={styles.kpiHeader}>
          <span className={styles.kpiTitle}>{title}</span>
          {isMet ? (
            <CheckCircle size={20} className={styles.kpiSuccess} />
          ) : (
            <XCircle size={20} className={styles.kpiPending} />
          )}
        </div>
        <div className={styles.kpiValues}>
          <span className={styles.kpiCurrent}>
            {format(current)}{unit}
          </span>
          <span className={styles.kpiTarget}>
            / {inverse ? '≤' : '≥'}{format(target)}{unit}
          </span>
        </div>
        <div className={styles.progressBar}>
          <div 
            className={`${styles.progressFill} ${isMet ? styles.progressSuccess : styles.progressPending}`}
            style={{ width: `${progress}%` }}
          />
        </div>
      </CardContent>
    </Card>
  );
}

export default async function KPIPage() {
  const { metrics } = await fetchData();

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>KPI Dashboard</h1>
          <p className={styles.subtitle}>
            12-month targets and progress tracking
          </p>
        </div>
        <Badge variant={metrics ? 'success' : 'warning'}>
          {metrics ? 'Live Data' : 'Demo Mode'}
        </Badge>
      </header>

      {/* Education KPIs */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <BookOpen size={24} />
          <h2 className={styles.sectionTitle}>Education</h2>
        </div>
        <div className={styles.kpiGrid}>
          <KPICard
            title="Workshops Delivered"
            current={metrics?.education.workshopsCount ?? 0}
            target={metrics?.education.workshopsTarget ?? 12}
          />
          <KPICard
            title="Participants Trained"
            current={metrics?.education.participantsTrained ?? 0}
            target={metrics?.education.participantsTarget ?? 300}
          />
          <KPICard
            title="Pre/Post Improvement"
            current={metrics?.education.prePostImprovement ?? 0}
            target={metrics?.education.improvementTarget ?? 25}
            unit="%"
            format={(v) => v.toFixed(1)}
          />
        </div>
      </section>

      {/* System KPIs */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <Clock size={24} />
          <h2 className={styles.sectionTitle}>System Performance</h2>
        </div>
        <div className={styles.kpiGrid}>
          <KPICard
            title="Median Report → Triage"
            current={metrics?.system.medianReportToTriage ?? 0}
            target={metrics?.system.triageTimeTarget ?? 30}
            unit=" min"
            inverse
            format={(v) => v.toFixed(1)}
          />
          <KPICard
            title="Verified Report Ratio"
            current={metrics?.system.verifiedRatio ?? 0}
            target={metrics?.system.verifiedRatioTarget ?? 60}
            unit="%"
            format={(v) => v.toFixed(1)}
          />
          <KPICard
            title="Abuse/False Report Rate"
            current={metrics?.system.abuseRate ?? 0}
            target={metrics?.system.abuseRateTarget ?? 5}
            unit="%"
            inverse
            format={(v) => v.toFixed(1)}
          />
          <KPICard
            title="Alert Publish Latency"
            current={metrics?.system.alertPublishLatency ?? 0}
            target={metrics?.system.publishLatencyTarget ?? 15}
            unit=" min"
            inverse
            format={(v) => v.toFixed(1)}
          />
        </div>
      </section>

      {/* Governance KPIs */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <Users size={24} />
          <h2 className={styles.sectionTitle}>Governance</h2>
        </div>
        <div className={styles.kpiGrid}>
          <KPICard
            title="Certified Triagers"
            current={metrics?.governance.certifiedTriagers ?? 0}
            target={metrics?.governance.triagersTarget ?? 15}
          />
        </div>
      </section>

      {/* Adoption KPIs */}
      <section className={styles.section}>
        <div className={styles.sectionHeader}>
          <Target size={24} />
          <h2 className={styles.sectionTitle}>Adoption</h2>
        </div>
        <div className={styles.kpiGrid}>
          <KPICard
            title="Partner Organizations"
            current={metrics?.adoption.partnerOrgs ?? 0}
            target={metrics?.adoption.partnerOrgsTarget ?? 4}
          />
          <KPICard
            title="External Deployments/Forks"
            current={metrics?.adoption.externalAdoption ?? 0}
            target={metrics?.adoption.externalAdoptionTarget ?? 2}
          />
        </div>
      </section>

      {/* Legend */}
      <div className={styles.legend}>
        <div className={styles.legendItem}>
          <CheckCircle size={16} className={styles.kpiSuccess} />
          <span>Target Met</span>
        </div>
        <div className={styles.legendItem}>
          <XCircle size={16} className={styles.kpiPending} />
          <span>In Progress</span>
        </div>
      </div>
    </div>
  );
}
