<?php

class SSE {
    public static function token($userId, $channel, $expires)
    {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, 'http://127.0.0.1:8001/token/');
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, [
            'userId' => $userId,
            'channel' => $channel,
            'exp' => $expires
        ]);
        $res = curl_exec($ch);
        curl_close($ch);
        return $res;
    }

    public static function emit($userId, $channel, $payload)
    {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, 'http://127.0.0.1:8001/emit/');
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, [
            'userId' => $userId,
            'channel' => $channel,
            'payload' => $payload
        ]);
        $res = curl_exec($ch);
        curl_close($ch);
    }
}
