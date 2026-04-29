import React from 'react';
import { InlineField, Input, Select, InlineSwitch } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource';
import { NovantQuery, NovantDataSourceOptions, QueryType } from '../types';

type Props = QueryEditorProps<DataSource, NovantQuery, NovantDataSourceOptions>;

const queryTypeOptions: Array<SelectableValue<QueryType>> = [
  { label: 'Trends', value: 'trends', description: 'Historical time series data' },
  { label: 'Live Values', value: 'values', description: 'Current point values' },
  { label: 'Points', value: 'points', description: 'Point metadata for a source or asset' },
  { label: 'Sources', value: 'sources', description: 'Source devices' },
  { label: 'Assets', value: 'assets', description: 'Equipment and assets' },
  { label: 'Spaces', value: 'spaces', description: 'Building spaces' },
  { label: 'Zones', value: 'zones', description: 'Building zones' },
];

const intervalOptions: Array<SelectableValue<string>> = [
  { label: 'Auto', value: 'auto' },
  { label: '5 min', value: '5min' },
  { label: '15 min', value: '15min' },
  { label: '30 min', value: '30min' },
  { label: '1 hour', value: '1hr' },
  { label: '1 day', value: '1day' },
  { label: '1 month', value: '1mo' },
  { label: 'Raw', value: 'raw' },
];

const aggregateOptions: Array<SelectableValue<string>> = [
  { label: 'Auto', value: 'auto' },
  { label: 'Mean', value: 'mean' },
  { label: 'Sum', value: 'sum' },
  { label: 'Min', value: 'min' },
  { label: 'Max', value: 'max' },
  { label: 'Diff', value: 'diff' },
];

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onQueryTypeChange = (val: SelectableValue<QueryType>) => {
    onChange({ ...query, queryType: val.value! });
    onRunQuery();
  };

  const onFieldChange = (field: keyof NovantQuery) => (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, [field]: event.target.value });
  };

  const onFieldBlur = () => {
    onRunQuery();
  };

  const onSelectChange = (field: keyof NovantQuery) => (val: SelectableValue<string>) => {
    onChange({ ...query, [field]: val.value });
    onRunQuery();
  };

  const onBoundOnlyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, boundOnly: event.target.checked });
    onRunQuery();
  };

  const { queryType } = query;

  return (
    <>
      <InlineField label="Query Type" labelWidth={14}>
        <Select
          options={queryTypeOptions}
          value={queryType}
          onChange={onQueryTypeChange}
          width={25}
        />
      </InlineField>

      {queryType === 'zones' && (
        <InlineField label="Zone IDs" labelWidth={14} tooltip="Comma-separated zone IDs (optional)">
          <Input
            value={query.zoneIds || ''}
            onChange={onFieldChange('zoneIds')}
            onBlur={onFieldBlur}
            placeholder="z.1,z.2"
            width={40}
          />
        </InlineField>
      )}

      {queryType === 'spaces' && (
        <InlineField label="Space IDs" labelWidth={14} tooltip="Comma-separated space IDs (optional)">
          <Input
            value={query.spaceIds || ''}
            onChange={onFieldChange('spaceIds')}
            onBlur={onFieldBlur}
            placeholder="sp.1,sp.2"
            width={40}
          />
        </InlineField>
      )}

      {queryType === 'assets' && (
        <InlineField label="Asset IDs" labelWidth={14} tooltip="Comma-separated asset IDs (optional)">
          <Input
            value={query.assetIds || ''}
            onChange={onFieldChange('assetIds')}
            onBlur={onFieldBlur}
            placeholder="a.1,a.2"
            width={40}
          />
        </InlineField>
      )}

      {queryType === 'sources' && (
        <>
          <InlineField label="Source IDs" labelWidth={14} tooltip="Comma-separated source IDs (optional)">
            <Input
              value={query.sourceIds || ''}
              onChange={onFieldChange('sourceIds')}
              onBlur={onFieldBlur}
              placeholder="s.1,s.2"
              width={40}
            />
          </InlineField>
          <InlineField label="Bound Only" labelWidth={14}>
            <InlineSwitch value={query.boundOnly || false} onChange={onBoundOnlyChange} />
          </InlineField>
        </>
      )}

      {(queryType === 'points' || queryType === 'values') && (
        <>
          <InlineField label="Source ID" labelWidth={14} tooltip="Source ID (or use Asset ID / Space ID)">
            <Input
              value={query.sourceId || ''}
              onChange={onFieldChange('sourceId')}
              onBlur={onFieldBlur}
              placeholder="s.1"
              width={25}
            />
          </InlineField>
          <InlineField label="Asset ID" labelWidth={14} tooltip="Asset ID (alternative to Source ID / Space ID)">
            <Input
              value={query.assetId || ''}
              onChange={onFieldChange('assetId')}
              onBlur={onFieldBlur}
              placeholder="a.1"
              width={25}
            />
          </InlineField>
          <InlineField label="Space ID" labelWidth={14} tooltip="Space ID (alternative to Source ID / Asset ID)">
            <Input
              value={query.spaceId || ''}
              onChange={onFieldChange('spaceId')}
              onBlur={onFieldBlur}
              placeholder="sp.1"
              width={25}
            />
          </InlineField>
          {queryType === 'values' && (
            <InlineField label="Point IDs" labelWidth={14} tooltip="Comma-separated point IDs (optional)">
              <Input
                value={query.pointIds || ''}
                onChange={onFieldChange('pointIds')}
                onBlur={onFieldBlur}
                placeholder="s.1.1,s.1.2"
                width={40}
              />
            </InlineField>
          )}
        </>
      )}

      {queryType === 'trends' && (
        <>
          <InlineField label="Point IDs" labelWidth={14} tooltip="Comma-separated point IDs (required)">
            <Input
              value={query.pointIds || ''}
              onChange={onFieldChange('pointIds')}
              onBlur={onFieldBlur}
              placeholder="s.1.1,s.1.2"
              width={40}
            />
          </InlineField>
          <InlineField label="Interval" labelWidth={14}>
            <Select
              options={intervalOptions}
              value={query.interval || 'auto'}
              onChange={onSelectChange('interval')}
              width={16}
            />
          </InlineField>
          <InlineField label="Aggregate" labelWidth={14}>
            <Select
              options={aggregateOptions}
              value={query.aggregate || 'auto'}
              onChange={onSelectChange('aggregate')}
              width={16}
            />
          </InlineField>
        </>
      )}
    </>
  );
}
