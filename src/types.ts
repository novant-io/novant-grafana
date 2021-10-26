//
// Copyright (c) 2021, Novant LLC
// Licensed under the MIT License
//
// History:
//   21 Oct 2021  Andy Frank  Creation
//

import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface NvQuery extends DataQuery {
  op: string;
  deviceId?: string;
  pointIds?: string;
}

export const defaultQuery: Partial<NvQuery> = {
  op: 'trends',
};

/**
 * These are options configured for each DataSource instance.
 */
export interface NvDataSourceOptions extends DataSourceJsonData {}

/**
 * Backend only secrets.
 */
export interface NvSecrets {
  apiKey?: string;
}
