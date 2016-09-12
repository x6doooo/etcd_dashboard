/**
 * Created by dx.yang on 16/9/12.
 */

(function () {


    var apis = {
        endpoints: '/api/endpoints/list'
    };


    angular.module('fe').service('endpointApiService', endpointApiService);

    /** @ngInject */
    function endpointApiService($ajax) {
        var me = this;
        me.getEndpoints = function() {
            return $ajax.get(apis.endpoints);
        };
    }

})();
