import React, { useEffect, useMemo, useState } from 'react';
import styled from 'styled-components';
import { ApiService } from '../../services/api';
import { FlagDebugItem } from '../../types/api';

const Container = styled.div`
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
`;

const PageHeader = styled.div`
  max-width: 600px;
  margin: 0 auto 1.5rem;
  color: white;
  text-align: center;
`;

const PageTitle = styled.h1`
  font-size: 1.6rem;
  margin: 0 0 0.5rem;
`;

const PageSubtitle = styled.div`
  font-size: 0.95rem;
  opacity: 0.9;
`;

const FlagsList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  max-width: 600px;
  margin: 0 auto;
`;

const FlagCard = styled.div`
  background: white;
  border-radius: 15px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  width: 100%;
`;

const CardHeader = styled.div`
  padding: 0.75rem 1rem;
  font-size: 0.95rem;
  color: #555;
  background: #fff;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const FlagContainer = styled.div`
  padding: 2rem;
  text-align: center;
  background: #f8f9fa;
`;

const QuestionText = styled.h2`
  color: #333;
  margin: 0 0 1.25rem 0;
  font-size: 1.5rem;
  text-align: center;
`;

const FLAG_MAX_WIDTH = 420;
const FLAG_MAX_HEIGHT = 260;

const FlagImage = styled.img`
  display: block;
  margin: 0 auto;
  max-width: ${FLAG_MAX_WIDTH}px;
  max-height: ${FLAG_MAX_HEIGHT}px;
  width: auto;
  height: auto;
`;

function svgToDataUrl(svg: string): string {
  const encoded = encodeURIComponent(svg)
    .replace(/'/g, '%27')
    .replace(/"/g, '%22');
  return `data:image/svg+xml;charset=utf-8,${encoded}`;
}

const StatusBox = styled.div`
  max-width: 600px;
  margin: 0 auto;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  text-align: center;
  color: #333;
`;

const FlagRow: React.FC<{ flag: FlagDebugItem }> = ({ flag }) => {
  const url = useMemo(
    () => (flag.flag_svg ? svgToDataUrl(flag.flag_svg) : ''),
    [flag.flag_svg]
  );
  return (
    <FlagCard>
      <CardHeader>
        <span>country_id: <b>{flag.country_id}</b></span>
        <span style={{ fontSize: '0.8rem', opacity: 0.7 }}>
          svg: {flag.flag_svg ? `${flag.flag_svg.length} chars` : 'empty'}
        </span>
      </CardHeader>
      <FlagContainer>
        <QuestionText>Guess The Flag</QuestionText>
        {url && <FlagImage src={url} alt={`Флаг ${flag.country_id}`} />}
      </FlagContainer>
    </FlagCard>
  );
};

const FlagsDebug: React.FC = () => {
  const [flags, setFlags] = useState<FlagDebugItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await ApiService.getAllFlags();
        if (!cancelled) setFlags(data);
      } catch (err: any) {
        if (!cancelled) setError(err?.response?.data?.error || err?.message || 'Не удалось загрузить флаги');
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <Container>
      <PageHeader>
        <PageTitle>Debug: все флаги</PageTitle>
        <PageSubtitle>
          Каждый флаг отрисован в том же контейнере, что и в игре. По одному на строку.
        </PageSubtitle>
      </PageHeader>

      {loading && <StatusBox>Загрузка флагов...</StatusBox>}
      {error && <StatusBox style={{ color: '#c33' }}>{error}</StatusBox>}

      {!loading && !error && (
        <>
          <PageHeader>
            <PageSubtitle>Всего флагов: {flags.length}</PageSubtitle>
          </PageHeader>
          <FlagsList>
            {flags.map((f) => (
              <FlagRow key={f.country_id} flag={f} />
            ))}
          </FlagsList>
        </>
      )}
    </Container>
  );
};

export default FlagsDebug;
