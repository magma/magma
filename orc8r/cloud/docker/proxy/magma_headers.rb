# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# Nghttpx can run mcruby scripts as part of each request handling:
# See: https://nghttp2.org/documentation/nghttpx.1.html?highlight=mruby-file#mruby-scripting

class App
  def on_req(env)
    # Inject Magma headers to inform the backend services about the client cert.
    # The headers would be present for all requests, and a empty string as
    # value indicate invalid cert.
    env.req.set_header("x-magma-client-cert-cn", env.tls_client_subject_name)
    env.req.set_header("x-magma-client-cert-serial", env.tls_client_serial.upcase)
  end
end

App.new
