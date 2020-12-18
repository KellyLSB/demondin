module.exports = {
  mode: 'development',
  entry: {
    checkout: `${__dirname}/jsx/checkout/main.jsx`,
    admin:    `${__dirname}/jsx/admin/main.jsx`,
  },
  output: { filename: '[name].js', path: `${__dirname}/public` },
  resolve: {
    extensions: [ '.wasm', '.mjs', '.ts', '.tsx', '.js', '.jsx', '.json' ]
  },
  devtool: 'source-map',
  module: {
    rules: [{
      test: /\.(t|j)s(|x)$/,
      use: {
        loader: 'babel-loader',
        options: {
          cacheDirectory: true,
          presets: [
            '@babel/preset-env',
            '@babel/preset-react',
            '@babel/preset-typescript',
            '@babel/preset-flow',
          ],
          plugins: [ '@babel/plugin-proposal-class-properties' ]
        }
      }
    },{
      test: /\.css$/,
      use: [
        { loader: "style-loader" },
        { loader: "css-loader" }
      ]
    },{
      test: /\.(png|jpg|gif|svg|woff|woff2|eot|ttf)$/i,
      use: [{
        loader: 'url-loader',
        options: { limit: 8192 }
      }]
    }]
  }
}
