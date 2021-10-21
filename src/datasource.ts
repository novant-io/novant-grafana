//
// Copyright (c) 2021, Novant LLC
// All Rights Reserved
//
// History:
//   21 Oct 2021  Andy Frank  Creation
//

import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { NvDataSourceOptions, NvQuery } from './types';

export class DataSource extends DataSourceWithBackend<NvQuery, NvDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NvDataSourceOptions>) {
    super(instanceSettings);
  }
}
