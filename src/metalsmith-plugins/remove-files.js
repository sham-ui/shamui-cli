const fs = require( 'fs' );
const path = require( 'path' );
const async = require( 'async' );

function removeFiles( files, metalsmith, done ) {
    const dest = metalsmith.metadata().dest;
    const keys = Object.keys( files );
    async.each( keys, function( file, next ) {
        delete files[ file ];
        fs.unlink( path.join( dest, file ), next );
    }, done );
}

module.exports = removeFiles;