//
// Copyright (c) 2023, Novant LLC
// Licensed under the MIT License
//
// History:
//   4 Dec 2023  Andy Frank  Creation
//

import React, { ChangeEvent } from 'react';
import { InlineField, Select, Input } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { NvDataSourceOptions, NvQuery } from '../types';

type Props = QueryEditorProps<DataSource, NvQuery, NvDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {

  // const onEndpointChange = (value: string) => {
  //   onChange({ ...query, endpoint: event.target.selectedValue });
  //};

  const onSourceIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, sourceId: event.target.value });
  };

  const onPointIdsChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, pointIds: event.target.value });
  };

  const { sourceId, pointIds } = query;
  const endpointOpts = [
    { label: 'Trends', value: 'trends' },
    // { label: 'Values', value: 'values' },
  ];

  return (
    <div className="gf-form">
      <InlineField label="Endpoint" labelWidth={14} tooltip="API endpoint for query">
        <Select onChange={() => console.log("todo")} options={endpointOpts} value="trends" />
      </InlineField>
      <InlineField label="Source ID" labelWidth={14} tooltip="Source ID to query">
        <Input onChange={onSourceIdChange} value={sourceId || ''} />
      </InlineField>
      <InlineField label="Point ID(s)" labelWidth={14} tooltip="Comma-separated list of point ID's to query">
        <Input onChange={onPointIdsChange} value={pointIds || ''} />
      </InlineField>
    </div>
  );
}
