<?php
	$answer_filename 	= dirname(__FILE__) . '/answer.dat';
	if (file_exists($answer_filename)) {
		$tmp_answer_filename= $answer_filename . '.tmp';
		copy($answer_filename, $tmp_answer_filename);
	}
	file_put_contents($answer_filename, "");

	$question_filename	= dirname(__FILE__) . '/question.txt';
	$fp_question		= fopen($question_filename, 'r');
	if (!$fp_question) {
		die;
	}
	
	$main_start_time = microtime(true);
	
	$fp_answer = @fopen($tmp_answer_filename, 'r');
	if ($fp_answer) {
		$question_index = 0;
		while (($question_line = fgets($fp_question)) !== false) {
			if ($question_index <= 1) {
				$question_index++;
				continue;
			}

			// 回答済みは飛ばす
			if (($answer_line = fgets($fp_answer)) !== false) {
				$answer_line = preg_replace("/\n|\r/", '', $answer_line);
				if ($answer_line != '') {
					$question_index++;
					file_put_contents($answer_filename, $answer_line . "\n", FILE_APPEND);
					continue;
				}
			}
	
			$question_line = preg_replace("/\n|\r/", '', $question_line);
			list($width, $height, $data) = explode(',', $question_line);

			$contents 	= '';
			$start_time = 0;
			$end_time	= 0;

			echo str_pad(($question_index - 1), 5, '0', STR_PAD_LEFT) . '問：';	
			
			$output = array();
			$start_time = microtime(true);
			exec(dirname(__FILE__) . '/6.out ' . $width . ' '. $height . ' ' . $data, $output);
			$end_time = microtime(true);
							
			if (array_key_exists(0, $output)) {
				$contents = $output[0];				
			}
					
			echo 'time:'	. ceil(($end_time - $start_time) * 1000) / 1000 . "\n";
			echo $contents . "\n";
				
			file_put_contents($answer_filename, $contents . "\n", FILE_APPEND);
			$question_index++;
		}
		
		@fclose($fp_answer);
	}
	
	@fclose($fp_question);
	
	echo 'time:'	. ceil((microtime(true) - $main_start_time) * 1000) / 1000 . "\n";
