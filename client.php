<?php 
$fp = stream_socket_client("tcp://127.0.0.1:9001", $errno, $errstr, 30);
if (!$fp) {
    echo "$errstr ($errno)<br />\n";
} else {
		$command = json_encode(array("cmd"=>"SET", "args"=>array(10000, 55)));
		$n = strlen($command);
		//tricky
		if(pack('L', 1) === pack('N', 1))
		{
			$bin = strrev(pack('l', $n));
		}else{
			$bin = pack('l', $n);
		}
		$bin .= $command;

    fwrite($fp, $bin);
    stream_set_blocking ($fp, 0);
    usleep(500000);
    $_responseHeader = @fread($fp, 4);
    if($_responseHeader)
    {
    	//tricky
			if(pack('L', 1) === pack('N', 1))
			{
				$length = unpack('llen', strrev($_responseHeader));
			}else{
				$length = unpack('llen', $_responseHeader);
			}

			$json = @fread($fp, $length['len']);
    }
	  
    
    stream_socket_shutdown($fp, STREAM_SHUT_WR);
    fclose($fp);
}

var_dump($json);

$jsonArray = json_decode($json, true);
$id = $jsonArray['Data'];

if(!is_numeric($id))
{
	var_dump("id error");exit;
}

// exit;

$fp = stream_socket_client("tcp://127.0.0.1:9001", $errno, $errstr, 30);
if (!$fp) {
    echo "$errstr ($errno)<br />\n";
} else {
		$command = json_encode(array("cmd"=>"GET", "args"=>array($id, "Julius")));
		$n = strlen($command);
		//tricky
		if(pack('L', 1) === pack('N', 1))
		{
			$bin = strrev(pack('l', $n));
		}else{
			$bin = pack('l', $n);
		}
		$bin .= $command;

    fwrite($fp, $bin);
    stream_set_blocking ($fp, 0);
    usleep(500000);
    $_responseHeader = @fread($fp, 4);
    if($_responseHeader)
    {
    	//tricky
			if(pack('L', 1) === pack('N', 1))
			{
				$length = unpack('llen', strrev($_responseHeader));
			}else{
				$length = unpack('llen', $_responseHeader);
			}

			$json = @fread($fp, $length['len']);
    }
	  
    
    stream_socket_shutdown($fp, STREAM_SHUT_WR);
    fclose($fp);
}

var_dump($json);

?>