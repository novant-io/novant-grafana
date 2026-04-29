import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { NovantQuery, NovantDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<NovantQuery, NovantDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NovantDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(): Partial<NovantQuery> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: NovantQuery, scopedVars: ScopedVars): NovantQuery {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      pointIds: query.pointIds ? templateSrv.replace(query.pointIds, scopedVars) : query.pointIds,
      sourceId: query.sourceId ? templateSrv.replace(query.sourceId, scopedVars) : query.sourceId,
      assetId: query.assetId ? templateSrv.replace(query.assetId, scopedVars) : query.assetId,
      spaceId: query.spaceId ? templateSrv.replace(query.spaceId, scopedVars) : query.spaceId,
      zoneIds: query.zoneIds ? templateSrv.replace(query.zoneIds, scopedVars) : query.zoneIds,
      spaceIds: query.spaceIds ? templateSrv.replace(query.spaceIds, scopedVars) : query.spaceIds,
      assetIds: query.assetIds ? templateSrv.replace(query.assetIds, scopedVars) : query.assetIds,
      sourceIds: query.sourceIds ? templateSrv.replace(query.sourceIds, scopedVars) : query.sourceIds,
    };
  }
}
