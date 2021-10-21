//
// Copyright (c) 2021, Novant LLC
// All Rights Reserved
//
// History:
//   21 Oct 2021  Andy Frank  Creation
//

import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface NvQuery extends DataQuery {
  deviceId?: string;
  pointIds?: string;
}

export const defaultQuery: Partial<NvQuery> = {};

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
