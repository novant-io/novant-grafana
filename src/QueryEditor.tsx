//
// Copyright (c) 2021, Novant LLC
// All Rights Reserved
//
// History:
//   21 Oct 2021  Andy Frank  Creation
//

import { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, NvDataSourceOptions, NvQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, NvQuery, NvDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onDeviceIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, deviceId: event.target.value });
  };

  onPointIdsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, pointIds: event.target.value });
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { deviceId, pointIds } = query;

    return (
      <div className="gf-form">
        <FormField
          labelWidth={8}
          value={deviceId || ''}
          onChange={this.onDeviceIdChange}
          label="Device ID"
          tooltip="Device to query"
        />
        <FormField
          labelWidth={8}
          value={pointIds || ''}
          onChange={this.onPointIdsChange}
          label="Point ID(s)"
          tooltip="Comma-separated list of point ID's to query"
        />
      </div>
    );
  }
}
