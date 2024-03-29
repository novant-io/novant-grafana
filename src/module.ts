//
// Copyright (c) 2023, Novant LLC
// Licensed under the MIT License
//
// History:
//   4 Dec 2023  Andy Frank  Creation
//

import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { NvQuery, NvDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, NvQuery, NvDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
