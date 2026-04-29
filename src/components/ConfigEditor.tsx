import React from 'react';
import { InlineField, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { NovantDataSourceOptions, NovantSecureJsonData } from '../types';

type Props = DataSourcePluginOptionsEditorProps<NovantDataSourceOptions, NovantSecureJsonData>;

export function ConfigEditor({ options, onOptionsChange }: Props) {
  const { secureJsonFields, secureJsonData } = options;

  const onAPIKeyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: { ...secureJsonData, apiKey: event.target.value },
    });
  };

  const onResetAPIKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: { ...secureJsonFields, apiKey: false },
      secureJsonData: { ...secureJsonData, apiKey: '' },
    });
  };

  return (
    <>
      <InlineField label="API Key" labelWidth={14} tooltip="Novant API key (starts with ak_)">
        <SecretInput
          isConfigured={Boolean(secureJsonFields?.apiKey)}
          value={secureJsonData?.apiKey || ''}
          placeholder="ak_..."
          width={40}
          onReset={onResetAPIKey}
          onChange={onAPIKeyChange}
        />
      </InlineField>
    </>
  );
}
