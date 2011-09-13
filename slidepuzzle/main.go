package main

import (
	"fmt"
	"strings"
	"sort"
	"math"
	"strconv"
	"flag"
	"container/heap"
	"container/vector"
	"time"
)

/**
 * GDD2011 チャレンジクイズ解答用プログラム
 * 
 * 以下のサイトを参考にさせていただきました。感謝！
 * http://www.ic-net.or.jp/home/takaken/nt/slide/solve15.html
 * http://www.geocities.jp/m_hiroi/xyzzy_lisp/abclisp14.html
 * http://d.hatena.ne.jp/g940425/20100812/1281624557
 * 
 * http://ja.wikipedia.org/wiki/%E5%8F%8D%E5%BE%A9%E6%B7%B1%E5%8C%96%E6%B7%B1%E3%81%95%E5%84%AA%E5%85%88%E6%8E%A2%E7%B4%A2
 * http://ja.wikipedia.org/wiki/%E6%B7%B1%E3%81%95%E5%84%AA%E5%85%88%E6%8E%A2%E7%B4%A2
 * 
 * 使い方
 * 6.out 3 3 123456780
 */
func main() {
	// コマンド引数の取得
	var has_error = false
	flag.Parse()
	
	// 幅の取得
	width, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		has_error = true
	}
	
	// 高さの取得
	height, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		has_error = true
	}

	// データの取得
	puzzle := flag.Arg(2);
	if len(puzzle) <= 0 {
		has_error = true
	}
	
	// データの取得
	is_debug, err := strconv.Atoi(flag.Arg(3));
	if err != nil {
		is_debug = 0;
	}
				
	// エラーがなければ、処理を開始する
	if !has_error {
		solver := new (Solver)
		solver.width 	= width
		solver.height	= height
		solver.puzzle	= puzzle
		solver.is_debug	= is_debug
	
		
	
		// 基本は上から順に詰めていくが、横幅のほうが大きい場合は左から詰めていく
		solver.to_top = true
		if solver.width > solver.height {
			solver.to_top = false
		}
		solver.is_repeat = false
		
		// 処理の開始時刻を取得
		solver.starttime = time.Nanoseconds()
		
		// 処理を開始する
		solver.Start()
	} else {	
		fmt.Println("error");	
	}
}

/**
 * 定数
 */
const (
	FIRST	int = 9
	
	UP		int = 0
	DOWN	int = 1
	LEFT	int = 2
	RIGHT	int = 3
	
	LEFT_TEXT	string = "L"
	RIGHT_TEXT	string = "R"
	UP_TEXT		string = "U"
	DOWN_TEXT	string = "D"
	FIRST_TEXT	string = "F"
	
	START int 			= 0
	WALL int 			= 100
	WALL_STRING string 	= "="
	
	MAX_TIME int64	= 0.2 * 1000 * 1000 * 1000
)

/**
 * パズル実行用構造体
 */
type Solver struct {
	starttime		int64				// 処理開始時間
	
	width			int					// 幅
	height			int					// 高さ
	puzzle			string				// 問題
	
	data			[]int				// 問題を配列に変換した物
	answer			[]int				// 解答を配列に変換した物
	md				[][][]int			// マンハッタン距離
	
	low				int					// 下限
	depth			int					// 深さ
	wall_count		int					// 壁の数
	wall			[][]bool			// 壁情報
	list			*vector.Vector		// 経路
	visited			map[string] bool	// 訪問チェック用
	search_count	int					// 探索回数
	start			*Cell				// 開始位置
	passage			int					// 通過状況
	to_top			bool				// 上・左探索状況用フラグ
	is_repeat		bool				// 再探索確認用フラグ
	
	is_debug		int					// デバッグ用フラグ
}


/**
 * 処理の実行
 */
func (solver *Solver) Start() {
	solver.initData()
	solver.initWall()
	
	solver.start = solver.createField(solver.data, FIRST, nil)
	var is_init bool
	solver.start.evaluation, is_init = solver.calcEvaluation(solver.start)
				
	if (is_init) {
		solver.Restart()		
	} else {
		// 枝切り用の設定
		solver.low 		= solver.start.evaluation
		solver.depth	= 1
	
		solver.list		= new(vector.Vector)
		heap.Init(solver.list)
		heap.Push(solver.list, solver.start)
	
		solver.passage			= 0
		solver.search_count 	= 0
		solver.calc()	
	}
}

/**
 * 処理の実行
 */
func (solver *Solver) Restart() {
	var is_init bool
	
	solver.start.evaluation, is_init = solver.calcEvaluation(solver.start)
	if (is_init) {
		solver.Restart()		
	} else {
		// 枝切り用の設定
		solver.low 		= solver.start.evaluation
		solver.depth	= 1
		
		solver.list		= new(vector.Vector)
		heap.Init(solver.list)
		heap.Push(solver.list, solver.start)
	
		solver.calc()
	}
}

/**
 * 計算処理
 */
func (solver *Solver) calc() {
	var index int
	var data []int
	var hash string
	var new_data []int
	var field *Cell
	var is_init bool
	
	solver.visited = make(map[string] bool, 1)
	for solver.list.Len() > 0 {
		cursor := heap.Pop(solver.list).(*Cell)
		solver.visited[cursor.hash] = true

		// 評価値が 0 （ゴール）に到着した
		if cursor.evaluation == 0 {
			solver.renderAnswer(cursor, solver.search_count)	
			return
		}
		
		index = cursor.space_position
		
		// 上 行けるか？
		if !solver.wall[index][UP] && cursor.direction != DOWN {
			data = cursor.data
			hash, new_data = solver.replace(data, index, index + (solver.width * -1))
			
			if !solver.visited[hash] {
				solver.visited[hash] = true
				field = solver.createField(new_data, UP, cursor.route)
				field.evaluation, is_init = solver.calcEvaluation(field)
				if (is_init) {					
					solver.Restart()
					return
				}
				if field.evaluation + field.depth <= solver.low + solver.depth {	
					heap.Push(solver.list, field)
				}
			}
		}
		
		// 下 行けるか？
		if !solver.wall[index][DOWN] && cursor.direction != UP {
			data = cursor.data
			hash, new_data = solver.replace(data, index, index + solver.width)
			if !solver.visited[hash] {
				solver.visited[hash] = true
				field = solver.createField(new_data, DOWN, cursor.route)
				field.evaluation, is_init = solver.calcEvaluation(field)
				if (is_init) {
					solver.Restart()
					return
				}
				if field.evaluation + field.depth <= solver.low + solver.depth {	
					heap.Push(solver.list, field)
				}
			}
		}
		
		// 左 行けるか？
		if !solver.wall[index][LEFT] && cursor.direction != RIGHT {
			data = cursor.data
			hash, new_data = solver.replace(data, index, index - 1)
			if !solver.visited[hash] {
				solver.visited[hash] = true
				field = solver.createField(new_data, LEFT, cursor.route)
				field.evaluation, is_init = solver.calcEvaluation(field)
				if (is_init) {
					solver.Restart()
					return
				}
				if field.evaluation + field.depth <= solver.low + solver.depth {	
					heap.Push(solver.list, field)
				}
			}
		}
		
		// 右 行けるか？
		if !solver.wall[index][RIGHT] && cursor.direction != LEFT {
			data = cursor.data
			hash, new_data = solver.replace(data, index, index + 1)
			if !solver.visited[hash] {
				solver.visited[hash] = true
				field = solver.createField(new_data, RIGHT, cursor.route)
				field.evaluation, is_init = solver.calcEvaluation(field)
				if (is_init) {
					solver.Restart()
					return
				}
				if field.evaluation + field.depth <= solver.low + solver.depth {	
					heap.Push(solver.list, field)
				}
			}
		}
		
		// 指定の解析時間を超えたら処理を中断する
		if (time.Nanoseconds() - solver.starttime > MAX_TIME) {
			// 上と左からどっちかのチェックが終わってなかったらもう一回チャレンジ
			if !solver.is_repeat {
				if solver.to_top {
					solver.to_top = false
				} else {
					solver.to_top = true
				}
				solver.is_repeat = true
				solver.starttime = time.Nanoseconds()
				solver.Start()
			}

			return
		}
	}
	
	solver.depth	+= 2
	solver.list		= new(vector.Vector)
	
	heap.Init(solver.list)
	heap.Push(solver.list, solver.start)
	
	solver.calc()
}

/**
 * フィールドの生成
 */
func (solver *Solver) createField(data []int, direction int, def_route []string) *Cell {
	var cell *Cell 		= new (Cell)
	cell.data 			= data
	
	// 経路保存用のデータ作成
	var route string = ""
	for i := 0; i < len(data); i++ {
		if data[i] >= WALL {
			route += "=:"
		} else {
			route += strconv.Itoa(int(data[i])) + ":"
		}
		
		if data[i] == 0 {
			cell.space_position = i
		}
	}		
	
	cell.hash			= route
	
	if direction == FIRST {
		route += FIRST_TEXT
	} else if direction == UP {
		route += UP_TEXT
	} else if direction == DOWN {
		route += DOWN_TEXT
	} else if direction == LEFT {
		route += LEFT_TEXT
	} else if direction == RIGHT {
		route += RIGHT_TEXT
	}
	
	// 経路の追加
	cell.route 			= def_route
	cell.AddRoute(route)
	cell.depth			= len(cell.route) - 1
	cell.direction		= direction
	
	return cell
}

/**
 * 評価値の算出
 */
func (solver *Solver) calcEvaluation(cell *Cell) (int, bool) {
	var start int
	var end int
	var distance int
	
	// 初期化フラグ
	is_init := false
	
	// 上から攻めるぞ
	if solver.to_top {		
		if (solver.passage < solver.height - 2) {
			start 	= solver.passage * solver.width;
			end		= (solver.passage + 1) * solver.width
			distance = 0
			for i := start; i < end	; i++ {
				if solver.answer[i] < WALL {
					for now_position := 0; now_position < len(solver.answer); now_position++ {
						if cell.data[now_position] == solver.answer[i] {
							distance += solver.md[solver.answer[i]][now_position][cell.space_position]
							break	
						}
					}
				}
			}
			
			if distance <= 0 {
				if solver.passage < solver.height - 2 {
					solver.passage++
					
					for x := 0; x < solver.width; x++ {
						position := (solver.passage * solver.width) + x
						solver.wall[position][UP] = true
					}

					solver.start = cell
					is_init = true
				}
			}
		} else {
			start 	= solver.passage * solver.width;
			end		= len(cell.data)
			distance = 0
			for i := start; i < end	; i++ {
				if cell.data[i] < WALL && cell.data[i] > 0 {
					distance += solver.md[cell.data[i]][i][cell.space_position]
				}
			}
		}
	// 左から攻めるぞ
	} else {
		if (solver.passage < solver.width - 2) {
			start 	= solver.passage
			end		= solver.height
			distance = 0
			
			for y := 0; y < solver.height; y++ {
				idx := (y * solver.width) + solver.passage
				if solver.answer[idx] < WALL {
					for now_position := 0; now_position < len(solver.answer); now_position++ {
						if cell.data[now_position] == solver.answer[idx] {
							distance += solver.md[solver.answer[idx]][now_position][cell.space_position]
							break	
						}
					}
				}
			}
			
			if distance <= 0 {				
				if solver.passage < solver.width - 2 {
					solver.passage++
					for y := 0; y < solver.height; y++ {
						position := (y * solver.width) + solver.passage
						solver.wall[position][LEFT] = true
					}
					
					solver.start = cell
					is_init = true
				}
			}
		} else {
			distance = 0
			for y := 0; y < solver.height; y++ {
				for x := solver.passage; x < solver.width; x++ {
					position := (y * solver.width) + x
					if cell.data[position] < WALL && cell.data[position] > 0 {
						distance += solver.md[cell.data[position]][position][cell.space_position]
					}
				}
			}
		}
	}
	return distance, is_init;
}

/**
 * パネルの置き換え
 */
func (solver *Solver) replace(data []int, oldPos int, newPos int) (string, []int) {
	new_data := make([]int, len(data))
	copy(new_data, data[0:])
	
	tmp := new_data[oldPos];
	new_data[oldPos] = new_data[newPos]
	new_data[newPos] = tmp 
	
	var route string = ""
	for i := 0; i < len(data); i++ {
		if new_data[i] >= WALL {
			route += "=:"
		} else {
			route += strconv.Itoa(new_data[i]) + ":"
		}
	}
	return route, new_data
}


/**
 * 回答の描画
 */
func (solver *Solver) renderAnswer(cursor *Cell, search_count int) {
	if solver.is_debug == 1 {
		fmt.Println("------------------------------------");
		fmt.Println(solver.data);
		fmt.Println(solver.answer);
		fmt.Println("------------------------------------");
	}
	
	for i := 0; i < len(cursor.route); i++ {
		text := strings.Split(cursor.route[i], "")
		if text[len(text) - 1] != FIRST_TEXT {
			fmt.Print(text[len(text) - 1]);
			if solver.is_debug == 1 {
				solver.renderMap(cursor.route[i])
			}
		} else {
			if solver.is_debug == 1 {
				solver.renderMap(cursor.route[i])
			}
		}
	}
	fmt.Print("\n");
}

/**
 * 回答の描画(デバッグ用)
 */
func (solver *Solver) renderMap(hash string) {
	fmt.Println("")
	data := strings.Split(hash, ":")
	index := 0;
	for height := 0; height < solver.height; height++ {
		for width := 0; width < solver.width; width++ {
			fmt.Print("[" + data[index] + "]")
			index++
		}
		fmt.Println("")
	}
}

/**
 * 進行方向の変換
 */
func (solver *Solver) convertDirection(direction int) string {
	direction_list := [...]string{"U","D","L","R"}
	return direction_list[direction]
}

/**
 * データの初期化
 */
func (solver *Solver) initData() {
	// 仮のテーブルを生成
	tmp := strings.Split(solver.puzzle, "")
	sort.Strings(tmp)
	wall := WALL

	solver.wall_count = 0
	var tmp2 map[string] int = make(map[string] int, len(tmp))
	var index int = 0
	for i := 0; i < len(tmp); i++ {
		if tmp[i] != WALL_STRING {
			tmp2[tmp[i]] = index
			index++
		} else {
			solver.wall_count++
		}
	}

	// 回答結果の生成
	solver.data 	= make([]int, len(solver.puzzle))
	solver.answer	= make([]int, len(solver.puzzle))
	
	tmpdata 		:= strings.Split(solver.puzzle, "")
	answer_idx 		:= 0
	index 			= 1
	
	for i := 0; i < len(tmpdata); i++ {
		if tmpdata[i] == "=" {
			solver.data[i] 				= wall
			solver.answer[answer_idx]	= wall
			wall++
		} else {
			solver.data[i] 				= tmp2[tmpdata[i]]
			solver.answer[answer_idx] 	= index
			index++
		}
		answer_idx++
	}
	solver.answer[len(solver.puzzle) - 1] = 0
	
	// マンハッタン距離の計算	
	solver.md		= make([][][]int, solver.height * solver.width)
	for i := 0;i < int(solver.height * solver.width); i++ {
		solver.md[i] = make([][]int, solver.height * solver.width)
		for j := 0;j < int(solver.height * solver.width); j++ {
			solver.md[i][j] = make([]int, solver.height * solver.width)	
		}
	}
	
	// 正解の位置と現在の位置と空白の位置を元に作成する
	for idx := 0; idx < len(solver.answer); idx++ {
		target := solver.answer[idx]
		
		if target > 0 && target < WALL {
			for height := 1; height <= solver.height; height++ {
				for width := 1; width <= solver.width; width++ {
					for space := 0; space < len(solver.puzzle); space++ {
						position := int(width) + (int(height) - 1) * int(solver.width)
	
						fwidth 	:= float64(solver.width)
						fpos 	:= float64(position - 1)
						fidx 	:= float64(idx)
						fspace	:= float64(space)
											
						distance_h := int(math.Fabs(math.Floor(fidx / fwidth) - math.Floor(fpos / fwidth)))
						if distance_h > 0 { 
							// 実際の移動距離に近づけてみる
							distance_h = distance_h * (distance_h + 1) / 2
						}
						distance_w := int(math.Fabs(float64(int(fidx) % int(fwidth) - int(fpos) % int(fwidth))))
						if distance_w > 0 {
							// 実際の移動距離に近づけてみる
							distance_w = distance_w * (distance_w + 1) / 2
						}
						
						space_w := 0
						space_h := 0
						if distance_h > 0 || distance_w > 0 {
							// 空白の場所も考慮してみる
							space_w	= int(math.Fabs(math.Floor(fspace / fwidth) - math.Floor(fidx / fwidth)))
							space_h = int(math.Fabs(float64(int(fspace) % int(fwidth) - int(fidx) % int(fwidth))))
						}
						solver.md[int(target)][int(position - 1)][space] = distance_h + distance_w + space_w + space_h
					}
				}
			}
		}
	}
}

/**
 * 壁の初期化
 */
func (solver *Solver) initWall() {
	// 縦・横のサイズの2次元配列を作成する
	solver.wall		= make([][]bool, solver.height * solver.width)
	for i := 0;i < solver.height * solver.width; i++ {
		// 2次元目には、上下左右の壁の状態を保存する
		solver.wall[i] = make([]bool, 4)
	}

	var y int
	var x int
	end := solver.width * solver.height - 1
	for y = 0; y < solver.height; y++ {
		for x = 0; x < solver.width; x++ {
			position := (y * solver.width) + x
			
			
			// 対象の位置が壁(=) なら上下左右すべて壁扱い
			if solver.data[position] >= WALL {
				solver.wall[position][UP]		= true
				solver.wall[position][DOWN]		= true
				solver.wall[position][LEFT]		= true
				solver.wall[position][RIGHT]	= true
			} else {
				solver.wall[position][UP]		= false
				solver.wall[position][DOWN]		= false
				solver.wall[position][LEFT]		= false
				solver.wall[position][RIGHT]	= false
				
				if y == 0 {
					solver.wall[position][UP]	= true
				} else if y == solver.height - 1 {
					solver.wall[position][DOWN]	= true
				}
				
				
				if x == 0 {
					solver.wall[position][LEFT]		= true
				} else if x == solver.width - 1 {
					solver.wall[position][RIGHT]	= true
				}
				
				position_up 		:= position - solver.width
				position_down 		:= position + solver.width
				position_left 		:= position - 1
				position_right		:= position + 1
				
				if position_up < 0 {
					solver.wall[position][UP]	= true
				} else if solver.data[position_up] >= WALL {
					solver.wall[position][UP]	= true
				}
				
				if position_down > end {
					solver.wall[position][DOWN]	= true
				} else if solver.data[position_down] >= WALL {
					solver.wall[position][DOWN]	= true
				}
				
				if position_left < 0 {
					solver.wall[position][LEFT]	= true
				} else if solver.data[position_left] >= WALL {
					solver.wall[position][LEFT]	= true
				}
				
				if position_right > end {
					solver.wall[position][RIGHT]	= true
				} else if solver.data[position_right] >= WALL {
					solver.wall[position][RIGHT]	= true
				}
			}
		}
	}
}

/**
 * セル構造体
 */
type Cell struct {
	data			[]int		// スライドパズルの状態
	depth			int			// 深さ
	direction		int			// 最終の進行方向
	route			[]string	// これまでの経路
	evaluation		int			// 評価
	space_position	int			// 空白マスの位置
	hash			string		// ハッシュ値
}

/**
 * 経路の追加
 */
func (cell *Cell) AddRoute(route string) {
	var tmp []string
	if cell.route == nil {
		// 経路の登録がない場合、新しい配列を生成する
		tmp 		= make([]string, 1)
		tmp[0] 		= route
		cell.route 	= tmp
	} else {
		// 経路の登録がある場合、現在の経路数 + 1の配列を生成する
		tmp = make([]string, len(cell.route) + 1)
		for i := 0; i < len(cell.route); i++ {
			tmp[i] = cell.route[i]
		}
		cell.route = tmp
		cell.route[len(cell.route) - 1] = route;
	}
}

/**
 * 初期位置の確認
 */
func (cell *Cell) CheckFirstPosition() {
	for i := 0; i < len(cell.data); i++ {
		if cell.data[i] == START {
			cell.space_position = i
			break
		}
	}
}

/**
 * ヒープ用の評価
 */
func (cell *Cell) Less(o interface{}) bool {
	// 深さ + 評価値 が 低い方を優先する
	return cell.evaluation + cell.depth < o.(*Cell).evaluation + o.(*Cell).depth
}
