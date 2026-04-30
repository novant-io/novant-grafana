import { DataQuery, DataSourceJsonData } from '@grafana/data';

export type QueryType = 'zones' | 'spaces' | 'assets' | 'sources' | 'points' | 'values' | 'trends';

export interface NovantQuery extends DataQuery {
  queryType: QueryType;
  // Entity filters
  zoneIds?: string;
  spaceIds?: string;
  assetIds?: string;
  sourceIds?: string;
  // Points/values context
  sourceId?: string;
  assetId?: string;
  spaceId?: string;
  pointIds?: string;
  pointTypes?: string;
  // Sources filter
  boundOnly?: boolean;
  // Trend options
  interval?: string;
  aggregate?: string;
}

export const DEFAULT_QUERY: Partial<NovantQuery> = {
  queryType: 'trends',
  interval: 'auto',
  aggregate: 'auto',
};

export interface NovantDataSourceOptions extends DataSourceJsonData {}

export interface NovantSecureJsonData {
  apiKey?: string;
}
