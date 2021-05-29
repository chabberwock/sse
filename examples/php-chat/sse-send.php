<?php
require __DIR__ . '/sse.php';
// This example is not secure, add validation of username via JWT or $_SESSION if you intend to use it in production
$username = htmlspecialchars($_POST['username']);
$message = htmlspecialchars($_POST['message']);

SSE::emit('*', 'chat', "event: message\ndata: " . json_encode(['username'=>$username, 'message' => $message]) . "\n");
