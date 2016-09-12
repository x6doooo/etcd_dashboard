/**
 * Created by dx.yang on 16/9/11.
 */


(function() {

    angular.module('fe')
        .controller('EndpointsListCtrl', EndpointsListCtrl);

    /** @ngInject */
    function EndpointsListCtrl($sce, endpointApiService) {

        var vm = this;


        endpointApiService.getEndpoints().done(function(data) {
            var list = data || [];
            var id2name = {};
            _.forEach(list, function(item) {
                id2name[item.Member.ID] = item.Member.name;
                item.peers = $sce.trustAsHtml(item.Member.peerURLs.join('<br>'));
                item.clients = $sce.trustAsHtml(item.Member.clientURLs.join('<br>'));
            });
            vm.list = list;
            vm.id2name = id2name;
        })

    }


})();