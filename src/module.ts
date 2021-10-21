//
// Copyright (c) 2021, Novant LLC
// All Rights Reserved
//
// History:
//   21 Oct 2021  Andy Frank  Creation
//

import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { NvQuery, NvDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, NvQuery, NvDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
