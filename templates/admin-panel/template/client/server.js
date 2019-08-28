const Express = require( 'express' );
const webpack = require( 'webpack' );
const webpackDevMiddleware = require( 'webpack-dev-middleware' );
const webpackHotMiddleware = require( 'webpack-hot-middleware' );
const config = require( './webpack.config' )( null, { mode: 'development' } );

const app = new Express();
const port = 3000;

const compiler = webpack( config );
app.use( webpackDevMiddleware( compiler, { noInfo: true, publicPath: config.output.publicPath } ) );
app.use( webpackHotMiddleware( compiler ) );

app.get( '/favicon.ico', function( req, res ) {
    res.sendFile( __dirname + '/favicon.ico' );
} );

app.get( '*', function( req, res ) {
    res.sendFile( __dirname + '/index.html' );
} );

app.listen( port, function( error ) {
    if ( error ) {
        console.error( error );
    } else {
        console.info(
            'Dev server for "{{ name }}" started.\n' +
            'Open up http://localhost:%s/ in your browser.',
            port
        );
    }
} );
