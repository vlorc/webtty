'use strict';

const path = require('path');
const webpack = require('webpack');

const BabiliWebpackPlugin = require('babili-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');

let config = {
    entry: {
        index: path.join(__dirname, 'lib/index.ts'),
    },
    module: {
        rules: [{
            test: /\.html$/,
            use: [{
                loader: "html-loader",
                options: { minimize: true }
            }]
        }, {
            enforce: 'pre',
            test: /\.js$/,
            loader: "source-map-loader",
        }, {
            test: /\.ts$/,
            use: 'ts-loader',
            exclude: /node_modules/,
        }, {
            test: /\.node$/,
            use: 'node-loader',
        }]
    },
    output: {
        filename: '[name].js',
        path: path.join(__dirname, 'dist'),
    },
    plugins: [
        new webpack.NoEmitOnErrorsPlugin(),
    ],
    resolve: {
        extensions: ['.js', '.jsx', '.json', '.ts', '.tsx'],
    },
    externals: [],
};

module.exports = function(env, argv) {
    config.plugins.push(
        new webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify(argv.mode),
        })
    );
    if (argv.mode === 'production') {
        config.plugins.push(
            new CleanWebpackPlugin(['./dist']),
            new BabiliWebpackPlugin(),
        );
    } else {
        config.devtool = 'source-map';
    }
    return config;
};