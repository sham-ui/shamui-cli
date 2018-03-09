const webpack = require( 'webpack' );

const plugins = [
    new webpack.NoEmitOnErrorsPlugin()
];

module.exports = {
    entry: {
        index: './src/{{name}}.js',
    },
    output: {
        path: __dirname,
        filename: '[name].js',
        publicPath: '/',
        library: ['{{name}}', '{{name}}/[name]' ],
        libraryTarget: 'umd'
    },
    externals: [
        'sham-ui'
    ],
    plugins: plugins,
    module: {
        loaders: [ {
            test: /(\.js)$/,
            loader: 'babel-loader',
            exclude: /(node_modules)/,
            include: __dirname
        }, {
            test: /\.sht/,
            loader: 'sham-ui-templates-loader'
        } ]
    }
};
