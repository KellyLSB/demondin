module.exports = {
  entry: __dirname + '/jsx/main.jsx',
  output: { filename: 'bundle.js', path: __dirname + '/public' },
  resolve: {
    extensions: [ '.wasm', '.mjs', '.js', '.jsx', '.json' ]
  },
  devtool: 'source-map',
  module: {
    rules: [{
      test: /\.js(|x)$/,
      use: {
        loader: 'babel-loader',
        options: {
          presets: [ '@babel/preset-env', '@babel/preset-react' ],
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