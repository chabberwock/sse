<?php
require __DIR__ . '/sse.php';
$u1 = ['Mighty','Shy','Heroic','Angry','Magic','Invisible','Strong','Handsome','Curious'];
$u2 = ['Panda','Koala','Cat','Penguin','Bear','Crow','Capuchin','Wolf'];
$username = $u1[array_rand($u1)] . $u2[array_rand($u2)] . rand(10,99);
$token = SSE::token($username, 'chat', 60*60*24);
?>

<script src="https://raw.githubusercontent.com/Yaffle/EventSource/master/src/eventsource.min.js"></script>
<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>

<script>
    $(document).ready(function() {
        var chatEvents = new EventSourcePolyfill("/api/sse/", {
            headers: {
                'Authorization': "Bearer <?php echo $token ?>"
            }
        });

        chatEvents.addEventListener('open', e => {
            console.log(e);
        });
        chatEvents.addEventListener('error', e => {
            console.log(e);
        });
        chatEvents.addEventListener('message', e => {
            console.log(e);
            var data = JSON.parse(e.data);
            $('#sse-chatbox').append('<div><span class="uname">'+data.username+'</span>: '+data.message+'</div>');
        });


        $('#sse-send').on('click', e => {
            e.preventDefault();
            $.post('/sse-send.php', {message: $('#message-input').val(), username: '<?php echo $username ?>'}, res => {
                $('#message-input').val('');
            });
        });
    });


</script>

