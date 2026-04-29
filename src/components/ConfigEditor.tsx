import React, { useState } from 'react';
import { Button, InlineField, SecretInput } from '@grafana/ui';
import { AppEvents, DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { getAppEvents, getBackendSrv } from '@grafana/runtime';
import { NovantDataSourceOptions, NovantSecureJsonData } from '../types';

type Props = DataSourcePluginOptionsEditorProps<NovantDataSourceOptions, NovantSecureJsonData>;

export function ConfigEditor({ options, onOptionsChange }: Props) {
  const { secureJsonFields, secureJsonData } = options;
  const [clearing, setClearing] = useState(false);

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

  const onClearCache = async () => {
    if (!options.uid) {
      return;
    }
    setClearing(true);
    try {
      await getBackendSrv().post(`/api/datasources/uid/${options.uid}/resources/clear-cache`);
      getAppEvents().publish({
        type: AppEvents.alertSuccess.name,
        payload: ['Point name cache cleared'],
      });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      getAppEvents().publish({
        type: AppEvents.alertError.name,
        payload: ['Failed to clear point cache', msg],
      });
    } finally {
      setClearing(false);
    }
  };

  return (
    <>
      <InlineField label="API Key" labelWidth={20} tooltip="Novant API key (starts with ak_)">
        <SecretInput
          isConfigured={Boolean(secureJsonFields?.apiKey)}
          value={secureJsonData?.apiKey || ''}
          placeholder="ak_..."
          width={40}
          onReset={onResetAPIKey}
          onChange={onAPIKeyChange}
        />
      </InlineField>
      <InlineField
        label="Point name cache"
        labelWidth={20}
        tooltip="Point names are cached per source for 24 hours to reduce API calls. Click to clear and force a refresh on the next query. Only takes effect after the data source has been saved."
      >
        <Button
          variant="secondary"
          size="sm"
          onClick={onClearCache}
          disabled={!options.uid || clearing}
        >
          {clearing ? 'Clearing…' : 'Clear cache'}
        </Button>
      </InlineField>
    </>
  );
}
