.. _rs_swagger:

Swagger UI
~~~~~~~~~~

The Digital Rebar Provision UI includes Swagger to allow exploration and testing of the :ref:`rs_api`.

To access the Swagger web interface, point your web browser at your DRP Endpoint IP addres, at the following URL:

  ::

    https://<IP_ADDRESS>:8092/swagger-ui

Ensure that the input form contains the same *IP_ADDRESS* reference as your DRP Endpoint.  Click the *Authorize* button to obtain an API Token from the Username/Password authentication.

See :ref:`rs_configuring_default` for default credentials.

.. note:: Due to the limitations of the Swagger specification, there are important API optimizations that are not exposed in the `swagger.json`.  Please review the :ref:`rs_api` for the complete capabilities.
