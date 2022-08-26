package main

import (
	"fmt"
	"regexp"

	"git.in.zhihu.com/zhsearch/search-ingress/pkg/common/ztext2"
)

const Template = `\<[\S\s]+?\>`

func main() {
	// ss:='1. ✐ ᵕ̈ ᵍᵒᵒᵈ ᴹᴼᴿᴺᴵᴺᴳ ᐝ ♡一路向阳 ๑҉ 今日份快乐正常营业 2. 早上好！偷我蚂蚁能量，揍我小鸡的人除外。 3. ᵍᵒᵒᵈ ᴹᴼᴿᴺᴵᴺᴳ ꔛꕤ❞ 任何你喜欢做的事，都不叫浪费时间。 4.每天早上醒来冲着空气挥一拳，今天要干翻这个世界。 5.︎ ɴɪᴄᴇ ᴅᴀʏ ︎ 别焦虑，总是会有不如意的， 生活从来都是泥沙俱下， 鲜花与荆棘并存， 我们带着诚意慢慢来！ 6. ᵕ̈ ᴹᴼᴿᴺᴵᴺᴳ 愿所有的早安都有回应 愿心中所愿都将被实现♡ ￼ 7. ♡生活再平凡” ʚ 也是限量版 ɞ 8. ◡̈ ᴴᴬᵛᴱ ᴬ ᴳᴼᴼᴰ ᵀᴵᴹᴱ 我们终会上岸，无论去到哪里 都是阳光万里，鲜花灿烂 ₍ᐢ •⌄• ᐢ₎ 9. ᵕ̈ ᴹᴼᴿᴺᴵᴺᴳ ᵕ̈ 在每一个充满希望的清晨， 告诉自己努力是为了遇见更好的自己，加油！ 10. “⛅️ ♡ 连起床这么艰难的事你都做到了， 接下来的一天还有什么能难倒你。 11. 慢生活，是有底气的自给自足，而不是好吃懒做的得过且过。早安！ 12. “起床才发现今天的阳光比昨天的更耀眼，就像你一样。” 13. 通往富婆之路的新一天开始了，早安~ 14. ᵕ̈ ᴹᴼᴿᴺᴵᴺᴳ ᵕ̈ 还不算晚的早餐 开启元气满满的一天嗷♡♡♡ 15. ♡good Morning 热爱世间万物，无最爱，无例外 16. 做一个积极向上的人，看温柔的句子，见阳光的人，眼里全是温柔和笑意。早安 ￼ 17. 起床就告诉你一个秘密：“今天的我依然喜欢你哦” 18. 一个人使劲踮起脚尖靠近太阳的时候，全世界都挡不住她的阳光。早安 19. 心里藏着小星星，生活才能亮晶晶。猫宁！ 20. ✧ '◡' ꀿᵖᵖᵞ 我来表演一个秒醒，早宀。 21. ☼ ⛅ ☼ gσσ∂ мσяиιиgꔛꕤ❞ 要努力呀，为了想要的生活， 为了人间的烟火，为了今天的风和月◞ᶫᵒᵛᵉ ♡ 22. ʜᴀ͟ᴘ͟ᴘ͟ʏ ᴇᴠᴇʀʏᴅᴀʏ̆̈ 今日，多云转甜 23. ️早吖～ 当努力到一定程度，幸运自会与你不期而遇。 24. ️早吖～ 不一定每天都很好，但每天都会有些小美好在等你~ 25. 如果生活踹了你好多脚，别忘了给它两个耳光，记住，奋斗总比流泪强。早安~️ 26. gσσ∂ мσяиιиg ya 人间烟火气，最抚凡人心。	'
	str := `<p></p><a href="https://link.zhihu.com/?target=https%3A//www.papertools.cn/app" data-draft-node="block" data-draft-type="link-card" class=" wrap external" target="_blank" rel="nofollow noreferrer">论文在线降重助手_论文降重工具 - 在线论文修改助手</a><p data-pid="6STrOzwD">上方是工具入口：</p><p data-pid="IdhbwhSd">论文查重检测报告结果出来后，打开把里面标红句子直接复制到上面，勾选3个转换选项，提交即可，转换出来的句子大概意思是不变的。</p><figure data-size="normal"><img src="https://pic3.zhimg.com/v2-6cbfc4e9e177a2c390a295df5ce197ae_b.jpg" data-caption="" data-size="normal" data-rawwidth="1150" data-rawheight="900" class="origin_image zh-lightbox-thumb" width="1150" data-original="https://pic3.zhimg.com/v2-6cbfc4e9e177a2c390a295df5ce197ae_r.jpg"/></figure><p></p><p></p>`
	// str := `<p data-pid="x8N0QxeK">1. ✐ ᵕ̈ ᵍᵒᵒᵈ ᴹᴼᴿᴺᴵᴺᴳ ᐝ</p><p><br/></p><p data-pid="RxT6Mq-k">♡一路向阳 ๑҉</p><p><br/></p><p data-pid="yIrokbEn">今日份快乐正常营业</p><p><br/></p><p data-pid="smicAG_w">2. 早上好！偷我蚂蚁能量，揍我小鸡的人除外。</p><p><br/></p><p data-pid="W09DsVYk">3. ᵍᵒᵒᵈ ᴹᴼᴿᴺᴵᴺᴳ ꔛꕤ❞</p><p><br/></p><p data-pid="TZyYPijP">任何你喜欢做的事，都不叫浪费时间。</p><p><br/></p><p data-pid="otiq7dwu">4.每天早上醒来冲着空气挥一拳，今天要干翻这个世界。</p><p><br/></p><p data-pid="3DcL4dQq">5.︎ ɴɪᴄᴇ ᴅᴀʏ ︎</p><p><br/></p><p data-pid="bz-UDL81">别焦虑，总是会有不如意的，</p><p><br/></p><p data-pid="PKcqa2E7">生活从来都是泥沙俱下，</p><p><br/></p><p data-pid="GX4-2a9E">鲜花与荆棘并存，</p><p><br/></p><p data-pid="SoXOV94s">我们带着诚意慢慢来！</p><p><br/></p><p data-pid="yRp-bpin">6. ᵕ̈ ᴹᴼᴿᴺᴵᴺᴳ</p><p><br/></p><p data-pid="lKCaikFc">愿所有的早安都有回应</p><p><br/></p><p data-pid="smi6X0e8">愿心中所愿都将被实现♡</p><p><br/></p><p data-pid="kLHpLybS">￼</p><p><br/></p><p data-pid="VgpfEVKF">7. ♡生活再平凡”</p><p><br/></p><p data-pid="5Jpys7SI">ʚ 也是限量版 ɞ</p><p><br/></p><p data-pid="MTTTfLrv">8. ◡̈ ᴴᴬᵛᴱ ᴬ ᴳᴼᴼᴰ ᵀᴵᴹᴱ</p><p><br/></p><p data-pid="fFvxRXCG">我们终会上岸，无论去到哪里</p><p><br/></p><p data-pid="VgVF6YNW">都是阳光万里，鲜花灿烂 ₍ᐢ •⌄• ᐢ₎</p><p><br/></p><p data-pid="LoHnIkAh">9. ᵕ̈ ᴹᴼᴿᴺᴵᴺᴳ   ᵕ̈  </p><p><br/></p><p data-pid="DsuPOYNz">在每一个充满希望的清晨，</p><p><br/></p><p data-pid="byujVKtN">告诉自己努力是为了遇见更好的自己，加油！</p><p><br/></p><p data-pid="b8-ODnir">10. “⛅️ ♡</p><p><br/></p><p data-pid="yVHHuZxg">连起床这么艰难的事你都做到了，</p><p><br/></p><p data-pid="GQ0pBPFH">接下来的一天还有什么能难倒你。</p><p><br/></p><p data-pid="0pR-TRmp">11. 慢生活，是有底气的自给自足，而不是好吃懒做的得过且过。早安！</p><p><br/></p><p data-pid="PJenICRy">12. “起床才发现今天的阳光比昨天的更耀眼，就像你一样。”</p><p><br/></p><p data-pid="4HQ2VjQ8">13. 通往富婆之路的新一天开始了，早安~</p><p><br/></p><p data-pid="-l_Exd20">14. ᵕ̈  ᴹᴼᴿᴺᴵᴺᴳ  ᵕ̈</p><p><br/></p><p data-pid="7gBKCSg7"> 还不算晚的早餐  </p><p><br/></p><p data-pid="_qbij4zA">开启元气满满的一天嗷♡♡♡</p><p><br/></p><p data-pid="tUxAnOvX">15. ♡good Morning</p><p><br/></p><p data-pid="XtpObCJ0">热爱世间万物，无最爱，无例外 </p><p><br/></p><p data-pid="Qs-WwKlT">16. 做一个积极向上的人，看温柔的句子，见阳光的人，眼里全是温柔和笑意。早安</p><p><br/></p><p data-pid="v6oRUNas">￼</p><p><br/></p><p data-pid="Yp9RvVRj">17. 起床就告诉你一个秘密：“今天的我依然喜欢你哦”</p><p><br/></p><p data-pid="f_tTk2Q0">18. 一个人使劲踮起脚尖靠近太阳的时候，全世界都挡不住她的阳光。早安</p><p><br/></p><p data-pid="nbtR3eEf">19. 心里藏着小星星，生活才能亮晶晶。猫宁！</p><p><br/></p><p data-pid="0g5X_aOw">20. ✧ &#39;◡&#39; ꀿªᵖᵖᵞ °</p><p><br/></p><p data-pid="KQDW2UCw">我来表演一个秒醒，早宀。</p><p><br/></p><p data-pid="CGIrKTLN">21. ☼ ⛅ ☼</p><p><br/></p><p data-pid="xJzm1Mub">gσσ∂ мσяиιиgꔛꕤ❞</p><p><br/></p><p data-pid="xK-e9QBp">要努力呀，为了想要的生活，</p><p><br/></p><p data-pid="TdYhiQ6M">为了人间的烟火，为了今天的风和月◞ᶫᵒᵛᵉ ♡</p><p><br/></p><p data-pid="yEEqi5FU">22. ʜᴀ͟ᴘ͟ᴘ͟ʏ ᴇᴠᴇʀʏᴅᴀʏ̆̈</p><p><br/></p><p data-pid="DrbRnv13">今日，多云转甜</p><p><br/></p><p data-pid="Id2LN4sl">23. ️早吖～</p><p><br/></p><p data-pid="es5mLa3J">当努力到一定程度，幸运自会与你不期而遇。</p><p><br/></p><p data-pid="adzRNI3i">24. ️早吖～</p><p><br/></p><p data-pid="uaTaeK7i">不一定每天都很好，但每天都会有些小美好在等你~</p><p><br/></p><p data-pid="vaOP6n4Y">25. 如果生活踹了你好多脚，别忘了给它两个耳光，记住，奋斗总比流泪强。早安~️</p><p><br/></p><p data-pid="cNJxYhsp">26. gσσ∂ мσяиιиg ya</p><p><br/></p><p data-pid="Kiw2Goac">人间烟火气，最抚凡人心。</p>`
	// str := "<p data-pid=\"rn31GDtY\">1 是非对错，有时只是角度问题。</p><p data-pid=\"uYmPM3FH\">2 越在意什么，什么越会折磨你。</p><p data-pid=\"2YMGDIx_\">3. 反感一个人，连听见名字都恶心。</p><p data-pid=\"oilJsb4f\"><b>4. 温和久了，稍有脾气就成了恶人。</b></p><p data-pid=\"wr_-z4jN\">5. 离去的都是风景，留下的才是人生。</p><img src=\"v2-a05a7e1d83cd246ab7885a5df7b13a06.png\" data-caption=\"\" data-size=\"normal\" data-rawwidth=\"639\" data-rawheight=\"480\" data-watermark=\"watermark\" data-original-src=\"v2-a05a7e1d83cd246ab7885a5df7b13a06\" data-watermark-src=\"v2-42970901fad86fb65c07894635fd4e67\" data-private-watermark-src=\"v2-5ee78ef55ef6f19e9748eb92b3f91f8e\"><p data-pid=\"6ayvs5Un\"><b>6. 任何安慰都没有自己看透来的奏效。</b></p><p data-pid=\"MKJVh5-P\">7. 我也想对每个人好，可狼我喂不饱。</p><p data-pid=\"fQitX0ob\">8. 友情都存在着吃醋，更别说爱情了。</p><p data-pid=\"rwNUJVrs\"><b>9. 原有前程可奔赴，亦有岁月可回首。</b></p><p data-pid=\"TrFUIBxE\">10. 偏见无非来自两个地方：无知和愚蠢。</p><img src=\"v2-857ce54830f6410fa1e8d356974cbbb2.png\" data-caption=\"\" data-size=\"normal\" data-rawwidth=\"638\" data-rawheight=\"479\" data-watermark=\"watermark\" data-original-src=\"v2-857ce54830f6410fa1e8d356974cbbb2\" data-watermark-src=\"v2-42811c5bd9d7c6afee9a9e262ae87eb8\" data-private-watermark-src=\"v2-0379881ccc5d39aa418dce651a1e913a\"><p data-pid=\"vV5uQQ2Q\">11. 其实一切的问题，时间已经给了答案。</p><p data-pid=\"UbEtdF0A\">12. 所有偷过的懒，都会变成打脸的巴掌。</p><p data-pid=\"X_QfeLhE\">13. 信任变得很难，假的和真的越来越像。</p><p data-pid=\"5dRltKB6\"><b>14. 理性的人适合共事，感性的人适合共处。</b></p><p data-pid=\"ZGt4pjMF\">15. 没有新故事的人，才会对过去念念不忘。</p><img src=\"v2-4226aa1418d635c54be3ace773fbb94a.png\" data-caption=\"\" data-size=\"normal\" data-rawwidth=\"638\" data-rawheight=\"479\" data-watermark=\"watermark\" data-original-src=\"v2-4226aa1418d635c54be3ace773fbb94a\" data-watermark-src=\"v2-bc79c851148d94f7ecbafa74fe736283\" data-private-watermark-src=\"v2-eefc32fc29869b7dc00f3619af2145fe\"><p data-pid=\"4HkQqllY\">16. 你嘴巴那么毒，内心一定是有很多苦吧。</p><p data-pid=\"7L1_P26w\"><b>17. 人生有两次很棒，虚惊一场和失而复得。</b></p><p data-pid=\"5fomldor\">18. 如果觉得生活苦，那就给自己撒点糖吧。</p><p data-pid=\"RDd1UwqT\">19. 如果没有见过光明，我本可以忍受黑暗。</p><p data-pid=\"RazPowtx\"><b>20. 谁心里没有故事，只不过是学会了控制。</b></p><img src=\"v2-56f97530ce469dc0b7ddec541af09713.png\" data-caption=\"\" data-size=\"normal\" data-rawwidth=\"638\" data-rawheight=\"479\" data-watermark=\"watermark\" data-original-src=\"v2-56f97530ce469dc0b7ddec541af09713\" data-watermark-src=\"v2-ea39c421a7ae8998bf2bdf5daddc4b96\" data-private-watermark-src=\"v2-07c3ff7886df052b3221cca3acb79ae8\"><p data-pid=\"4bpq0RNK\">21. 我觉得单身挺好，可是一直却招人讽刺。</p><p data-pid=\"YzMRxaGG\">22. 我欣赏你的直言不讳，但请管好你的嘴。</p><p data-pid=\"_HytS-5s\"><b>23. 希望曾经的仰望，都能成为以后的日常。</b></p><p data-pid=\"_tYOLkei\">24. 心若没有栖息的地方，到哪里都是流浪。</p><p data-pid=\"G9CpbyZW\">25. 越长大你越知道，有钱比有什么都舒坦。</p><img src=\"v2-0fc3d5052d74abec5e7bce64c5a6da83.png\" data-caption=\"\" data-size=\"normal\" data-rawwidth=\"638\" data-rawheight=\"479\" data-watermark=\"watermark\" data-original-src=\"v2-0fc3d5052d74abec5e7bce64c5a6da83\" data-watermark-src=\"v2-0d0f7573de3817cf2ea4371105392f9a\" data-private-watermark-src=\"v2-039bbae0cd8c86726a366f3c5841f6be\"><p data-pid=\"WvITciK2\">26. 两个人的沟通，70%是情绪，30%是内容。</p><p data-pid=\"OYIVoDjR\"><b>27. 不发生点烂事，永远看不清身边人的模样。</b></p><p data-pid=\"PWxOdyzp\">28.对于我们最爱的人，不说永远，只说珍惜。</p><p data-pid=\"4CCMhLVk\">29. 多少次的义无反顾，在时间面前也得认输。</p><p data-pid=\"GUeQpk_t\">30. 我会等，因为最好的东西，总是压轴出场。</p><img src=\"v2-ff3137deafb84f40abda9d13d8db1c07.png\" data-caption=\"\" data-size=\"normal\" data-rawwidth=\"638\" data-rawheight=\"479\" data-watermark=\"watermark\" data-original-src=\"v2-ff3137deafb84f40abda9d13d8db1c07\" data-watermark-src=\"v2-d37602302471f88820256b0b7a7f0d44\" data-private-watermark-src=\"v2-1d513c91589d02d39c2a94d2554da80d\"><p data-pid=\"sQoPfFN8\">31. 幸福是比较级，要有东西垫底才感觉的到。</p><p data-pid=\"57amR7Gj\">32. 在泥泞的道路上，要保持一颗玩泥巴的心。</p><p data-pid=\"MF48mHDT\"><b>33. 这世界上唯一扛得住岁月摧残的就是才华。</b></p>"
	reg := regexp.MustCompile("<p.*?>\\d[.．,、 ]+(.*?)</p.*?>")
	reg1 := regexp.MustCompile("<br/>(.*?)<br/>")
	// reg2 := regexp.MustCompile("<br/>(.*?)<br/>")
	r := regexp.MustCompile(Template)
	// fmt.Println(r.ReplaceAllString(str, "\n"))
	for _, str := range reg.FindAllStringSubmatch(str, -1) {
		if len(str) == 2 {
			fmt.Println(r.ReplaceAllString(str[1], ""))
		}
	}

	fmt.Println("=============================================")

	for _, str := range reg.FindAllStringSubmatch(str, -1) {
		if len(str) == 2 {
			fmt.Println(ztext2.ZTextPlaintext(str[1]))
		}
	}
	fmt.Println("-------------------------------------------")

	for _, str := range reg.FindAllStringSubmatch(str, -1) {
		if len(str) == 2 {
			fmt.Println(str[1])
			fmt.Println(reg1.FindAllStringSubmatch(str[1], -1))
		}
	}
}
