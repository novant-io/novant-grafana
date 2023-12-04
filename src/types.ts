//
// Copyright (c) 2023, Novant LLC
// Licensed under the MIT License
//
// History:
//   4 Dec 2023  Andy Frank  Creation
//

import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface NvQuery extends DataQuery {
  endpoint: string;
  sourceId?: string;
  pointIds?: string;
}

export const DEFAULT_QUERY: Partial<NvQuery> = {
  endpoint: "trends"
};

/**
 * These are options configured for each DataSource instance
 */
export interface NvDataSourceOptions extends DataSourceJsonData {}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface NvSecrets {
  apiKey?: string;
}
