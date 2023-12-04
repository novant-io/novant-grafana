//
// Copyright (c) 2023, Novant LLC
// Licensed under the MIT License
//
// History:
//   4 Dec 2023  Andy Frank  Creation
//

import { DataSourceInstanceSettings, CoreApp } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { NvQuery, NvDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<NvQuery, NvDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NvDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<NvQuery> {
    return DEFAULT_QUERY;
  }
}
