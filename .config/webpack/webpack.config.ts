import type { Configuration } from 'webpack';
import path from 'path';
import CopyWebpackPlugin from 'copy-webpack-plugin';
import ForkTsCheckerWebpackPlugin from 'fork-ts-checker-webpack-plugin';

const config = async (env: Record<string, unknown>): Promise<Configuration> => {
  const isProduction = Boolean(env.production);

  return {
    cache: {
      type: 'filesystem',
      buildDependencies: { config: [path.resolve(process.cwd(), '.config/webpack/webpack.config.ts')] },
    },
    context: path.join(process.cwd(), 'src'),
    devtool: isProduction ? 'source-map' : 'eval-source-map',
    entry: { module: './module.ts' },
    externals: [
      'lodash', 'jquery', 'moment', 'react', 'react-dom', 'react-redux', 'redux', 'rxjs',
      '@grafana/data', '@grafana/ui', '@grafana/runtime', '@grafana/e2e-selectors',
      '@emotion/react', '@emotion/css',
    ],
    mode: isProduction ? 'production' : 'development',
    module: {
      rules: [
        {
          exclude: /node_modules/,
          test: /\.[tj]sx?$/,
          use: {
            loader: 'swc-loader',
            options: {
              jsc: {
                baseUrl: path.resolve(process.cwd()),
                target: 'es2021',
                loose: false,
                parser: { syntax: 'typescript', tsx: true, decorators: false, dynamicImport: true },
              },
            },
          },
        },
        { test: /\.css$/, use: ['style-loader', 'css-loader'] },
        {
          test: /\.(png|jpe?g|gif|svg)$/,
          type: 'asset/resource',
          generator: {
            filename: isProduction ? '[hash][ext]' : '[file]',
            publicPath: 'public/plugins/novant-datasource/',
          },
        },
      ],
    },
    output: {
      clean: { keep: /gpx_|plugin.json/ },
      filename: '[name].js',
      library: { type: 'amd' },
      path: path.resolve(process.cwd(), 'dist'),
      publicPath: '/',
      uniqueName: 'novant-datasource',
    },
    plugins: [
      new CopyWebpackPlugin({
        patterns: [
          { from: 'plugin.json', to: '.' },
          {
            from: 'img',
            to: 'img',
            noErrorOnMissing: true,
            globOptions: { ignore: ['**/.DS_Store'] },
          },
        ],
      }),
      new ForkTsCheckerWebpackPlugin({
        async: !isProduction,
        issue: { include: [{ file: '**/*.{ts,tsx}' }] },
        typescript: { configFile: path.join(process.cwd(), 'tsconfig.json') },
      }),
    ],
    resolve: {
      extensions: ['.js', '.jsx', '.ts', '.tsx'],
      unsafeCache: true,
    },
  };
};

export default config;
