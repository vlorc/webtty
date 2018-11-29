'use strict'

const path = require('path');
const webpack = require('webpack');

const BabiliWebpackPlugin = require('babili-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');

let config = {
    entry: {
        main: path.join(__dirname, 'src/main.ts'),
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
        }, {
            test: /\.(gif|jpg|png|woff|svg|eot|ttf)\??.*$/,
            loader: 'url-loader?limit=8192&name=img/[name].[ext]'
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
}

module.exports = function(env, argv) {
    config.module.rules.push({
        test: /\.css$/,
        use: [argv.mode === 'production' ? MiniCssExtractPlugin.loader : 'style-loader'],
    });

    config.plugins.push(
        new webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify(argv.mode),
        })
    );
    if (argv.mode === 'production') {
        config.plugins.push(
            new CleanWebpackPlugin(['./dist']),
            new BabiliWebpackPlugin(),
            new MiniCssExtractPlugin({
                filename: "[name].css",
                chunkFilename: "[id].css",
            }),
            new HtmlWebpackPlugin({ template: path.resolve(__dirname, 'index.html') }),
        );
    } else {
        config.devtool = 'source-map';
    }
    config.plugins.push(
        new HtmlWebpackPlugin({ template: path.resolve(__dirname, 'index.html') }),
    );
    return config;
};