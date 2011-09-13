<?php
	$question_filename	= dirname(__FILE__) . '/question.txt';
	$fp_question		= fopen($question_filename, 'r');
	if (!$fp_question) {
		die;
	}
	
	$main_start_time = microtime(true);
	
	$answer_count	= 0;
	$success_count 	= 0;
	$fail_count		= 0;
	$answer_filename 	= dirname(__FILE__) . '/answer.dat';
	$fp_answer = @fopen($answer_filename, 'r');
	if ($fp_answer) {
		// 2行飛ばす
		$question_line = fgets($fp_question);
		$question_line = fgets($fp_question);
		
		$count = 1;
		while (($question_line = fgets($fp_question)) !== false) {
			$question_line = preg_replace("/\n|\r/", '', $question_line);
			list($width, $height, $data) = explode(',', $question_line);

			echo "問" . $count . "\n";
			$count++;

			$tmp 	= array();
			$puzzle = array();
			$start	= 0;
			for ($i = 0; $i < strlen($data); $i++) {
				$puzzle[] = substr($data, $i, 1);
				if (substr($data, $i, 1) === '0') {
					$start = $i;
				} else if (substr($data, $i, 1) !== '=') {
					$tmp[] = substr($data, $i, 1);
				}
			}
			sort($tmp);
			
			$idx = 0;
			$answer = '';
			for ($i = 0; $i < strlen($data); $i++) {
				if (substr($data, $i, 1) === '=') {
					$answer .= '=';
				} else {
					if ($idx < count($tmp)) {
						$answer .= $tmp[$idx];
						$idx++;
					}
				}
			}
			$answer .= '0';
			
			// 回答済みは飛ばす
			$pos = $start;
			if (($answer_line = fgets($fp_answer)) !== false) {
				$answer_line = preg_replace("/\n|\r/", '', $answer_line);
				if ($answer_line != '') {
					for ($i = 0; $i < strlen($answer_line); $i++) {
						$direction = substr($answer_line, $i, 1);
						if ($direction == 'U') {
							$tmp = $puzzle[$pos];
							$puzzle[$pos] = $puzzle[$pos - $width];
							$puzzle[$pos - $width] = $tmp;
							$pos = $pos - $width;
						} else if ($direction == 'D') {
							$tmp = $puzzle[$pos];
							$puzzle[$pos] = $puzzle[$pos + $width];
							$puzzle[$pos + $width] = $tmp;						
							$pos = $pos + $width;
						} else if ($direction == 'L') {
							$tmp = $puzzle[$pos];
							$puzzle[$pos] = $puzzle[$pos - 1];
							$puzzle[$pos - 1] = $tmp;
							$pos = $pos - 1;
						} else if ($direction == 'R') {
							$tmp = $puzzle[$pos];
							$puzzle[$pos] = $puzzle[$pos + 1];
							$puzzle[$pos + 1] = $tmp;
							$pos = $pos + 1;
						}
					}

					if (implode('', $puzzle) == $answer) {
						echo '正解' . "\n";
						$success_count++;
					} else {
						echo $data . "\n";
						echo $answer_line . "\n";
						echo implode('', $puzzle) . "\n";
						echo $answer . "\n";
						echo '不正解' . "\n";
						$fail_count++;
					}
					$answer_count++;
				}
			}
		}
		
		@fclose($fp_answer);
	}
	
	@fclose($fp_question);
	
	echo $success_count . ':' . $fail_count . ':' . $answer_count;
