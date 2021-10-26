//
// Copyright (c) 2021, Novant LLC
// Licensed under the MIT License
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

const { Select, FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, NvQuery, NvDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  // onOpChange = (event: ChangeEvent<HTMLSelectElement>) => {
  //   const { onChange, query, onRunQuery } = this.props;
  //   onChange({ ...query, op: event.currentTarget.value });
  //   onRunQuery();
  // };

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
    const opOptions = [
      { label: 'Trends', value: 'trends' },
      // { label: 'Values', value: 'values' },
    ];
    const query = defaults(this.props.query, defaultQuery);
    const { op, deviceId, pointIds } = query;
    return (
      <div className="gf-form">
        <label className="gf-form-label width-3">Op</label>
        <Select
          options={opOptions}
          value={opOptions[op === 'values' ? 1 : 0]}
          // onChange={this.onOpChange}
          // tooltip="API operand for query"
        />
        <FormField
          labelWidth={6}
          value={deviceId || ''}
          onChange={this.onDeviceIdChange}
          label="Device ID"
          tooltip="Device to query"
        />
        <FormField
          labelWidth={7}
          value={pointIds || ''}
          onChange={this.onPointIdsChange}
          label="Point ID(s)"
          tooltip="Comma-separated list of point ID's to query"
        />
      </div>
    );
  }
}
