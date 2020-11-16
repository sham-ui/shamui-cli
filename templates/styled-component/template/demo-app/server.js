const Express = require( 'express' );
const webpack = require( 'webpack' );
const webpackDevMiddleware = require( 'webpack-dev-middleware' );
const webpackHotMiddleware = require( 'webpack-hot-middleware' );
const config = require( './webpack.config' );

const app = new Express();
const port = 3000;

const compiler = webpack( config );
app.use( webpackDevMiddleware( compiler, { noInfo: true, publicPath: config.output.publicPath } ) );
app.use( webpackHotMiddleware( compiler ) );

app.get( '*', function( req, res ) {
    res.sendFile( __dirname + '/index.html' );
} );

app.listen( port, function( error ) {
    if ( error ) {

        // eslint-disable-next-line no-console
        console.error( error );
    } else {

        // eslint-disable-next-line no-console
        console.info(
            'App "Demo for {{name}}" started.\n' +
            'Open up http://localhost:%s/ in your browser.',
            port
        );
    }
} );
