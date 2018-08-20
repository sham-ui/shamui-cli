const path = require( 'path' );
const webpack = require( 'webpack' );
const ExtractTextPlugin = require( 'extract-text-webpack-plugin' );

const plugins = [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin(),
    new ExtractTextPlugin( {
        allChunks: true,
        filename: 'bundle.css'
    } )
];
const entry = [
    path.join( __dirname, 'src/main.js' ),
    path.join( __dirname, 'src/styles/main.scss' ),
    'webpack-hot-middleware/client'
];

module.exports = {
    entry,
    mode: 'development',
    devtool: 'cheap-module-eval-source-map',
    output: {
        path: path.join( __dirname, 'dist' ),
        filename: 'bundle.js',
        publicPath: '/dist/'
    },
    plugins: plugins,
    module: {
        rules: [ {
            test: /\.scss$/,
            loader: ExtractTextPlugin.extract( {
                fallback: 'style-loader',
                use: 'css-loader!sass-loader'
            } )
        }, {
            test: /\.css$/,
            loader: ExtractTextPlugin.extract( {
                fallback: 'style-loader',
                use: 'css-loader'
            } )
        }, {
            test: /(\.js)$/,
            loader: 'babel-loader',
            exclude: /(node_modules)/,
            include: [ __dirname, path.join( __dirname, '../src' ) ]
        }, {
            test: /\.sht$/,
            loader: 'sham-ui-templates-loader?hot'
        } ]
    }
};
