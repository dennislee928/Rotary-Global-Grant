import { GraduationCap, Users, TrendingUp, Target, Calendar, MapPin } from 'lucide-react';
import { getTrainingEvents, getTrainingStats } from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import styles from './page.module.css';

async function fetchData() {
  try {
    const [stats, events] = await Promise.all([
      getTrainingStats().catch(() => null),
      getTrainingEvents({ pageSize: 10 }).catch(() => null),
    ]);
    return { stats, events: events?.data || [] };
  } catch {
    return { stats: null, events: [] };
  }
}

export default async function TrainingPage() {
  const { stats, events } = await fetchData();

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>Training Hub</h1>
          <p className={styles.subtitle}>
            Anti-fraud education and community safety workshops
          </p>
        </div>
      </header>

      {/* Stats Overview */}
      <div className={styles.statsGrid}>
        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <Calendar size={24} />
              </div>
              <div className={styles.statInfo}>
                <span className={styles.statValue}>{stats?.totalEvents ?? '--'}</span>
                <span className={styles.statLabel}>Workshops</span>
                <span className={styles.statTarget}>Target: 12</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <Users size={24} />
              </div>
              <div className={styles.statInfo}>
                <span className={styles.statValue}>{stats?.totalParticipants ?? '--'}</span>
                <span className={styles.statLabel}>Participants</span>
                <span className={styles.statTarget}>Target: 300</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <TrendingUp size={24} />
              </div>
              <div className={styles.statInfo}>
                <span className={styles.statValue}>
                  {stats?.averageImprovement ? `+${stats.averageImprovement.toFixed(1)}%` : '--'}
                </span>
                <span className={styles.statLabel}>Avg Improvement</span>
                <span className={styles.statTarget}>Target: +25%</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card variant="elevated">
          <CardContent>
            <div className={styles.statCard}>
              <div className={styles.statIcon}>
                <Target size={24} />
              </div>
              <div className={styles.statInfo}>
                <Badge variant={stats?.targetMet ? 'success' : 'warning'}>
                  {stats?.targetMet ? 'On Track' : 'In Progress'}
                </Badge>
                <span className={styles.statLabel}>KPI Status</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Training Modules */}
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Training Modules</h2>
        <div className={styles.modulesGrid}>
          <Card className={styles.moduleCard}>
            <CardContent>
              <div className={styles.moduleIcon}>
                <GraduationCap size={32} />
              </div>
              <h3 className={styles.moduleTitle}>Recognizing Phishing</h3>
              <p className={styles.moduleDescription}>
                Learn to identify phishing emails, fake websites, and social engineering tactics.
              </p>
              <div className={styles.moduleMeta}>
                <Badge size="sm">30 min</Badge>
                <Badge size="sm" variant="info">Beginner</Badge>
              </div>
            </CardContent>
          </Card>

          <Card className={styles.moduleCard}>
            <CardContent>
              <div className={styles.moduleIcon}>
                <GraduationCap size={32} />
              </div>
              <h3 className={styles.moduleTitle}>Safe Reporting</h3>
              <p className={styles.moduleDescription}>
                How to report incidents safely without escalating panic or spreading misinformation.
              </p>
              <div className={styles.moduleMeta}>
                <Badge size="sm">25 min</Badge>
                <Badge size="sm" variant="info">Beginner</Badge>
              </div>
            </CardContent>
          </Card>

          <Card className={styles.moduleCard}>
            <CardContent>
              <div className={styles.moduleIcon}>
                <GraduationCap size={32} />
              </div>
              <h3 className={styles.moduleTitle}>Crisis Communication</h3>
              <p className={styles.moduleDescription}>
                Best practices for communicating during emergencies and high-stress situations.
              </p>
              <div className={styles.moduleMeta}>
                <Badge size="sm">45 min</Badge>
                <Badge size="sm" variant="warning">Intermediate</Badge>
              </div>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Recent Events */}
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Recent Workshops</h2>
        {events.length > 0 ? (
          <div className={styles.eventsGrid}>
            {events.map((event) => (
              <Card key={event.id}>
                <CardContent>
                  <div className={styles.eventHeader}>
                    <h3 className={styles.eventTitle}>{event.title}</h3>
                    {event.improvement && (
                      <Badge variant="success">+{event.improvement.toFixed(1)}%</Badge>
                    )}
                  </div>
                  <div className={styles.eventMeta}>
                    <span className={styles.eventMetaItem}>
                      <Calendar size={14} />
                      {event.eventDate}
                    </span>
                    <span className={styles.eventMetaItem}>
                      <MapPin size={14} />
                      {event.location}
                    </span>
                    <span className={styles.eventMetaItem}>
                      <Users size={14} />
                      {event.attendanceCount} attendees
                    </span>
                  </div>
                  {event.preAvg !== undefined && event.postAvg !== undefined && (
                    <div className={styles.eventScores}>
                      <div className={styles.scoreItem}>
                        <span className={styles.scoreLabel}>Pre-test</span>
                        <span className={styles.scoreValue}>{event.preAvg.toFixed(1)}%</span>
                      </div>
                      <div className={styles.scoreArrow}>→</div>
                      <div className={styles.scoreItem}>
                        <span className={styles.scoreLabel}>Post-test</span>
                        <span className={styles.scoreValue}>{event.postAvg.toFixed(1)}%</span>
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <Card>
            <CardContent>
              <p className={styles.emptyState}>No training events recorded yet</p>
            </CardContent>
          </Card>
        )}
      </section>
    </div>
  );
}
