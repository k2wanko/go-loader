const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const isProd = process.env.NODE_ENV === 'production'

module.exports = {
    entry: path.join(__dirname, 'app.go'),
    output: {
        filename: 'bundle.js'
    },
    module: {
        loaders: [
            {
                test: /\.go$/,
                exclude: /node_modules/,
                loader: 'go-loader'
            }
        ]
    },
    devtool: isProd ? '#source-map' : '#eval',
    plugins: [new HtmlWebpackPlugin()]
}

if (isProd) {
    module.exports.plugins = (module.exports.plugins || []).concat([
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: '"production"'
            }
        }),
        new webpack.optimize.UglifyJsPlugin({
            sourceMap: true,
            compress: {
                warnings: false
            }
        }),
        new webpack.LoaderOptionsPlugin({
            minimize: true
        })
    ])
}