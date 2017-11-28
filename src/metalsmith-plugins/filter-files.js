const evaluate = require( '../utils/eval' );

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

module.exports = filterFiles;