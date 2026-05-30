-module(logger_ffi).
-export([add_file_handler/1]).

add_file_handler(Filename) ->
    logger:add_handler(file_logger, logger_std_h, #{
        config => #{
            type => file,
            file => binary_to_list(Filename)
        }
    }),
    nil.
