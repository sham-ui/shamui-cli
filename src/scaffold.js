const path = require( 'path' );
//const minimatch = require( 'minimatch' );
const Metalsmith = require( 'metalsmith' );
const async = require( 'async' );
const render = require( 'consolidate' ).handlebars.render;
const evaluate = require( './eval' );
const getOptions = require( './options' );
const prompt = require( './prompt' );

function scaffold( name, src, dest, done ) {
    const opts = getOptions( name, src );
    const metalsmith = Metalsmith( path.join( src, 'template' ) );

    const data = Object.assign( metalsmith.metadata(), {
        destDirName: name,
        inPlace: dest === process.cwd(),
        noEscape: true
    } );

    metalsmith.use( promptQuestions( opts.prompts ) )
        .use( filterFiles( opts.filters ) )
        .use( renderTemplates( opts.skipInterpolation ) );

    metalsmith.clean( false )
        .source( '.' ) // start from template root instead of `./src` which is Metalsmith's default for `source`
        .destination( dest )
        .build( done );

    return data;
}

function promptQuestions( prompts ) {
    return function( files, metalsmith, done ) {
        prompt( prompts, metalsmith.metadata(), done );
    }
}

function filter( files, filters, data, done ) {
    if ( !filters ) {
        return done();
    }

    const fileNames = Object.keys( files );

    Object.keys( filters ).forEach( function( glob ) {
        fileNames.forEach( function( file ) {
            if ( match( file, glob, { dot: true } ) ) {
                const condition = filters[ glob ];
                if ( !evaluate( condition, data ) ) {
                    delete files[ file ];
                }
            }
        } );
    } );

    done();
}

function filterFiles( filters ) {
    return function( files, metalsmith, done ) {
        filter( files, filters, metalsmith.metadata(), done )
    }
}

function renderTemplates( skipInterpolation ) {
    if ( typeof skipInterpolation === 'string' ) {
        skipInterpolation = [ skipInterpolation ];
    }

    return function( files, metalsmith, done ) {
        const keys = Object.keys( files );
        const metalsmithMetadata = metalsmith.metadata();

        async.each( keys, function( file, next ) {
            if ( skipInterpolation &&
                multimatch( [ file ], skipInterpolation, { dot: true } ).length ) {
                return next();
            }

            const str = files[ file ].contents.toString();

            if ( !/{{([^{}]+)}}/g.test( str ) ) {
                return next();
            }

            render( str, metalsmithMetadata, function( err, res ) {
                if ( err ) {
                    err.message = `[${file}] ${err.message}`;
                    return next( err );
                }
                files[ file ].contents = new Buffer( res );
                next();
            } );

        }, done );
    }
}

module.exports = scaffold;